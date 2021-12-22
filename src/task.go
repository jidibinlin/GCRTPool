package GCRTPool

type task struct {
	function func(params ...interface{})
	params   []interface{}
	release  bool
}

// pack ...
// func NewTask(function func(params []interface{}), params ...interface{}) *task {
// 	foo := new(task)
// 	foo.function = function
// 	foo.params = params
// 	return foo
// }

func CoRun(function func(params ...interface{}), params ...interface{}) {
	foo := new(task)
	foo.function = function
	foo.params = params
	GetMgr().Process(foo)
}
