package logger


// 日志输出接口
type LogOuter interface{

}

// 控制台输出器
type ConsoleLogOuter struct{

}

type FileLogOuter struct{

}
// 时间切分文件输出器
type TimeCutFileLogOuter struct{

}
// 容量切分文件输出器
type CapacityCutFileLogOuter struct{

}
