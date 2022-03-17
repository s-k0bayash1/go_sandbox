package approximation_element

func init() {
	hoge(1)
	//hoge(MyInt(1)) can't compile
	fuga(1)
	fuga(MyInt(1))
	hogehoge(1) // can compile
	hogehoge(MyInt(1))
	//hogehoge(MyInt2(1)) can't compile
	//hogehoge(MyInt3(1)) can't compile
}

type NotArrowUnderlyingTypeInt interface {
	int
}

type ArrowUnderlyingTypeInt interface {
	~int
}

type NotArrowUnderlyingTypeMyInt interface {
	MyInt
}

// can't compile
//type ArrowUnderlyingTypeMyInt interface {
//	~MyInt
//}

type MyInt int
type MyInt2 int
type MyInt3 MyInt

func hoge[T NotArrowUnderlyingTypeInt](_ T) {}

func fuga[T ArrowUnderlyingTypeInt](_ T) {}

func hogehoge[T NotArrowUnderlyingTypeMyInt](_ T) {}
