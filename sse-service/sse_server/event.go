package sse_server

type SseEvent struct {
    Broadcast  bool
    ReceiverId string
    Data       []byte
}

func NewEvent(receiverId string, data []byte) *SseEvent {
    return &SseEvent{
        Broadcast:false,
        ReceiverId:receiverId,
        Data:data,
    }
}

func NewBroadcastEvent(data []byte) *SseEvent {
    return &SseEvent{
        Broadcast:true,
        ReceiverId:"",
        Data:data,
    }
}