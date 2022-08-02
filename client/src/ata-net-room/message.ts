
export enum ResponseType {
    GameEvtRes = 0,
    JoinReject = 1,
    YouJoinRes = 2,
    PlayerJoinRes = 3,
    LeaveRes = 4,
    LockRes  = 5,
    UnlockRes = 6,
};

export enum RequestType {
    GameEvtReq = 0,
    LeaveReq = 1,
    LockReq = 2,
    UnlockReq = 3
};

export interface WSResponse {
    type: ResponseType
    send: number
    relay: number
    sender: number
    id?: string
    name?: string
    payload?: string
    parsedPayload: any
};

export interface WSRequest {
    type: RequestType
    send: number
    id?: string
    name?: string
    payload?: string
    parsedPayload: any
};

export const parseWSResponse = (data: string): WSResponse|undefined => {
    const res = JSON.parse(data) as WSResponse;
    if(res?.payload) {
      res.parsedPayload = JSON.parse(atob(res.payload));
    }

    return res
}

export const responseTypeName = (type: ResponseType) => {
    return Object.keys(ResponseType).filter(k => isNaN(Number(k)))[type];
}