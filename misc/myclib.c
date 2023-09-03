#include "stdio.h"

extern long long my_sum(long long a, long long b);


void print_str(char *s)
{
    int t = my_sum(1,2);
    printf("%s\n", s?s:"nil");
}