package main

var CONCURRENT_TASK_NUM = 5

var CONCURRENT_LIST_PARSER_NUM = 5

var CONCURRENT_DETAIL_PAGE_HANDLER_NUM = 5

var TASK_CAPACITY = CONCURRENT_TASK_NUM * CONCURRENT_LIST_PARSER_NUM * CONCURRENT_DETAIL_PAGE_HANDLER_NUM

func crawlEngineConcurrentTaskCapacity() int {
	return TASK_CAPACITY
}
