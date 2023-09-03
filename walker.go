package main

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"modules/transform_html"
	"modules/utils"

	"golang.org/x/sync/errgroup"
)

var py_sub_rx *regexp.Regexp = regexp.MustCompile(`\\\d+`)
var py_sub_rx_plcmt = func(s string) string {
	offset := 0
	x, _ := strconv.Atoi(s[1:])
	x = x + offset
	return "$" + strconv.Itoa(x)
}

var NotFoundError = errors.New("Source not found")

type Walker struct {
	source_name string
	etl_config_path string
	etl_config *EtlConfig
	source_config *SourceConfig
	request_maker *utils.RequestMaker
	menu_page_url_sub string
}

func (w *Walker) Init(source_name string, etl_config_path string) (*Walker, error) {
	var err error
	w.source_name = source_name
	w.etl_config_path = etl_config_path

	w.etl_config, err = w.parse_config(etl_config_path)
	if err != nil { 
		return nil, err
	}

	w.source_config, err = w.extract_source_config(source_name)
	if err != nil { 
		return nil, err
	}
	hc := w.etl_config.Http
	w.request_maker = (&utils.RequestMaker{}).Init(
		&utils.InitParams{
			Timeout: time.Second * time.Duration(hc.Retries.Timeout),
			StatusForcelist: hc.Retries.status_forcelist,
			BackoffFactor: hc.Retries.Backoff_factor,
			MaxRetries: hc.Retries.Max_retries,
			Headers: hc.Headers,
		},
	)
	w.menu_page_url_sub = py_sub_rx.ReplaceAllStringFunc(w.source_config.Menu.Page_url_sub, py_sub_rx_plcmt)
	return w, nil
}

func (w *Walker) parse_config(location_path string) (*EtlConfig, error) {
	file, err := ioutil.ReadFile(location_path)
	if err != nil {
		return nil, err
	}
	var config EtlConfig
	err = yaml.Unmarshal(file, &config)
	return &config, err
}

func (w *Walker) extract_source_config(source_name string) (source_config *SourceConfig, err error) {
	if err != nil {
		return	
	}
	for _, sc := range w.etl_config.Sources {
		if sc.Name == source_name {
			source_config = &sc
			return
		}
	}
	err = NotFoundError
	return
}

func (w* Walker) extract_data(url string, rule []*transform_html.ParserTransfromRule) (ret_data *transform_html.TransformMap, err error) {
	resp, err := w.request_maker.Request(&utils.RequestParams{
		Method: `GET`,
		Url: url,
		BodyReader: nil,
	})
	ret_data = &map[string]any{}
	err = transform_html.TransformHtmlReader(ret_data, resp.Body, rule)
	defer resp.Body.Close()
	return
}

func (w *Walker) sub_page_number(filter_url string, page_number uint) string {
	if page_number < 2 {
		return utils.UrlCombineAll(w.source_config.Root_url, filter_url, w.source_config.Menu.First_page_url)
	}
	page_url := regexp.MustCompile(`(\d+)`).ReplaceAllString(strconv.Itoa(int(page_number)), w.menu_page_url_sub)
	return utils.UrlCombineAll(w.source_config.Root_url, filter_url, page_url)
}

func (w *Walker) parse_menu_page(filter_url string, page_number uint) (*transform_html.TransformMap, error) {
	url := w.sub_page_number(filter_url, page_number)
	return w.extract_data(url, w.source_config.Menu.Rules)
}

func (w *Walker) parse_card_page(url_part string) (*transform_html.TransformMap, error) {
	return w.extract_data(utils.UrlCombine(w.source_config.Root_url, url_part), w.source_config.Card.Rules)
}

func extract_value[T any](m *transform_html.TransformMap, key string) T {
	if value, ok := (*m)[key]; ok {
		return value.(T)
	}
	var value T
	return value
}

type ConsumerType = func (data []*transform_html.TransformMap, page_num uint)

func (w *Walker) WalkSync(filter_url string, begin, end uint, consumer ConsumerType) (err error) {
	for num := begin; num < (end +1); num++ {
		num := num
		err = walk_on_menu_page(w, filter_url, num, consumer)
		if err != nil {
			return
		}
	}
	return
}


func (w *Walker) Walk(filter_url string, begin, end uint, consumer ConsumerType) (err error) {
	grp := new(errgroup.Group)
	mtx := sync.Mutex{}
	for num := begin; num < (end +1); num++ {
		num := num
		grp.Go( func() error {
			return walk_on_menu_page(w, filter_url, num, func(data []*transform_html.TransformMap, page_num uint) {
				mtx.Lock()
				defer mtx.Unlock()
				consumer(data, page_num)
			})
		})
	}
	err = grp.Wait()
	return
}


func walk_on_menu_page(w *Walker, filter_url string, num uint, consumer ConsumerType) (err error) {
	var menu *transform_html.TransformMap
	menu, err = w.parse_menu_page(filter_url, num)
	if err != nil {
		return err
	}
	menu_items := extract_value[*transform_html.TransformList](menu, `menu_items`)
	if menu_items == nil {
		consumer([]*map[string]any{}, num)
		return nil
	}
	card_data_list := []*transform_html.TransformMap{}
	for _, card_any := range *menu_items {
		card := card_any.(*transform_html.TransformMap)
		var card_data *transform_html.TransformMap
		card_data, err = w.parse_card_page((*card)[`url`].(string))
		if err != nil {
			return err
		}
		card_data_list = append(card_data_list, card_data)
	}
	consumer(card_data_list, num)
	return nil
}