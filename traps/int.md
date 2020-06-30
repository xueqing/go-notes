# int 类型在 64 位机器是 8 字节

在 64 位机器上，int 是 8 字节。陷阱示例：结构体RTPMuxer 保存的 ssrc 是 uint32，需要转成 int 传给 ffmpeg，但是当 ssrc 最高位是 1 时直接传给 ffmpeg 会溢出，超出 INT_MAX。这是因为golang 的 uint32 直接转成 int 使用了超过 4 字节保存，除了原来的 4 字节，还额外使用一个符号位。需要先转成 int32，确定使用 4 字节有效位保存数据。最后 ffmpeg 只使用这 4 个有效字节转成 uint32。
