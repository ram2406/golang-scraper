package transform_html

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/html"
)

func _py_adopt_rx(rx [2]string) [2]string {
	plcmnt := py_sub_rx.ReplaceAllStringFunc(rx[1], py_sub_rx_plcmt)
	return [2]string { rx[0], plcmnt }
}
func _py_adopt_selector(sel string) string {
	return strings.ReplaceAll(sel, `:-soup-contains`, `:contains`)
}

func _handle_regex(rxp [2]string, text string) string {
	rx := _py_adopt_rx(rxp)
	rx_ := regexp.MustCompile(rx[0])
	return rx_.ReplaceAllString(text, rx[1])
}

func _handle_attr(selected_soup goquery.Selection, attr_name string) string {
	res, exists := selected_soup.Attr(attr_name)
	if !exists {
		return ""
	}
	return res
}

func _transform_html_1(
	transformed_data IListMapUnion,
	soup *html.Node,
	rule []*ParserTransfromRule,
	level int,
	limit int,
) ( err error ) {

	for _, r := range rule {
		err = _transform_html_2(transformed_data, soup, r, level, limit)
		if err != nil {
			return
		}
	}
	return
}



func _transform_html_2(
	transformed_data IListMapUnion,
	soup *html.Node,
	rule *ParserTransfromRule,
	level int,
	limit int,
) ( err error ) {
	if level > limit {
		err = errors.New(fmt.Sprintf(`level [%d] greater than limit [%d]`, level, limit))
		return
	}
	
	selected_soup := goquery.NewDocumentFromNode(soup)
	var transformed_data_out IListMapUnion = transformed_data

	if rule.Grouping != `` && !transformed_data.is_list() {
		transformed_data_out = transformed_data
	}

	if rule.Grouping != `` && transformed_data.is_list() {
		transformed_data_out = ListWrapper{}
		transformed_data.pack(transformed_data_out.extract(), ``)
	}
	
	if rule.Selector != `` {
		selected := selected_soup.Find(_py_adopt_selector(rule.Selector))
		tags := selected.Nodes
		if len(tags) == 0 && rule.Exception_on_not_found == false {
			html, _ := selected_soup.Html()
			log.Debugf("\n not found %s", rule.Selector)
			log.Debugf("\n %s", html[:50])
			return
		}
		if len(tags) == 0 && rule.Exception_on_not_found == true {
			msg := fmt.Sprintf(`None found at all by selector "%s"`, rule.Selector)
			html, _ := selected_soup.Html()
			log.Debugf(`%s \n %s \n %#v`, msg, html, rule)
			err = errors.New(msg)
		}
		
		rule_override := *rule
		rule_override.Selector = ``
			
		if len(tags) > 1 && transformed_data.is_list() {
			for _, tag := range tags {
				err = _transform_html_2(transformed_data_out, tag, &rule_override, level +1, limit)
				if err != nil {
					return
				}
			}
			return
		}
		if len(tags) > 1 {
			nested_data := &ListWrapper{}
		
			for _, tag := range tags {
				err = _transform_html_2(nested_data, tag, &rule_override, level +1, limit)
				if err != nil {
					return
				}
			}
			key_name := rule.Mapping
			if rule.Grouping != `` {
				key_name = rule.Grouping
			}
			transformed_data_out.pack(nested_data.extract(), key_name)
			return
		}
		selected_soup = goquery.NewDocumentFromNode(tags[0])
	}

	key_name := rule.Mapping

	if key_name == `` && len(rule.Children) == 0 {
		key_name = rule.Grouping
	}

	if key_name != `` {
		attr_name := rule.Attribute_Name
		if attr_name == `` {
			attr_name = `text`
		}
		text := ``
		if attr_name == `text` {
			text = selected_soup.Text()
		} else {
			text = _handle_attr(*selected_soup.Selection, attr_name)
		}
		text = strings.TrimSpace(text)
		handled_text := text
		if rule.Regex_Sub_Value[0] != `` {
			handled_text = _handle_regex(rule.Regex_Sub_Value, text)
		}

		transformed_data_out.pack(handled_text, key_name)
	}

	if len(rule.Children) > 1 {
		for _, r := range rule.Children {
			err = _transform_html_2(transformed_data_out, selected_soup.Nodes[0], r, level +1, limit)
			if err != nil {
				return
			}
		}
		
	}

	return
}

func TransformHtml(
	transformed_data *map[string]any,
	soup *goquery.Document,
	rule []*ParserTransfromRule,
) (error) {
	data := MapWrapper{ mapw: transformed_data }
	return _transform_html_1(data, soup.Get(0), rule, 1, 1000)
}


func TransformHtmlText(
	transformed_data *map[string]any,
	soup string,
	rule []*ParserTransfromRule,
) (error) {
	return TransformHtmlReader(transformed_data, strings.NewReader(soup), rule)
}


func TransformHtmlReader(
	transformed_data *map[string]any,
	soup io.Reader,
	rule []*ParserTransfromRule,
) (error) {
	doc, err := goquery.NewDocumentFromReader(soup)
	if err != nil {
		return err
	}
	return TransformHtml(transformed_data, doc, rule,)
}



func TransformHtmlList(
	transformed_data *[]any,
	soup *goquery.Document,
	rule []*ParserTransfromRule,
) {
	data := ListWrapper{listw: transformed_data}
	_transform_html_1(data, soup.Get(0), rule, 1, 1000)
}
