import React, { useRef, useState } from 'react';
import './App.scss';
import { parseWSResponse, RequestType, ResponseType, responseTypeName, WSRequest, WSResponse } from './ata-net-room/message';
import { v4 as uuidv4 } from 'uuid';

function App() {
  const serverRef = useRef<HTMLInputElement>(null);
  const codeRef = useRef<HTMLInputElement>(null);
  const createRef = useRef<HTMLInputElement>(null);
  const geInputRef = useRef<HTMLTextAreaElement>(null);
  const nameRef = useRef<HTMLInputElement>(null);

  const encoder = new TextEncoder();
  const decoder = new TextDecoder();

  const [client, setClient] = useState<WebSocket | undefined>(undefined);
  const [responses, setResponses] = useState<WSResponse[]>([]);
  const [roomLocked, setRoomLocked] = useState<boolean>(false);

  const renderResponses = (responses: WSResponse[]) => {
    return <div className='flex-column'>{responses.map(r => (
      <div key={r.id} className="card">
        <div>ID: {r.id}</div>
        <div>Sender: {r.sender}</div>
        <div>Sent at: {new Date(r.send).toTimeString()}</div>
        <div>Relayed at: {new Date(r.relay).toTimeString()}</div>
        <div>Type: {responseTypeName(r.type)} ({r.type})</div>
        <div>Payload: {JSON.stringify(r.parsedPayload)}</div>
      </div>))}
    </div>;
  };

  const connected = () => client !== undefined;

  const sendRequest = (ws: WebSocket, dataObj: WSRequest | any) => {
    const json = JSON.stringify(dataObj);
    console.log(`sending json: ${json}`);
    const bytes = encoder.encode(json);
    ws.send(bytes)
  }

  const initClient = () => {
    const ws = new WebSocket(serverRef!.current!.value, "atanet_v1");
    ws.onopen = function () {
      const joinReq = {
        room_code: codeRef.current!.value,
        create: createRef.current!.checked,
        player_name: nameRef.current!.value,
        game_id: "React Tester",
        room_size: 16
      };

      sendRequest(this, joinReq);
    }
    ws.onclose = () => {
      setClient(undefined);
      console.log("closed!")
      setRoomLocked(false);
    }

    ws.onmessage = (msg) => {
      const res = parseWSResponse(msg.data);
      if (!res) {
        console.error("failed to parse ws response");
        return;
      }

      setResponses(oldRes => [res, ...oldRes]);

      console.log(`Got response: ${responseTypeName(res.type)}`);

      switch (res.type) {
        case ResponseType.LeaveRes:
          client?.close();
          break
        case ResponseType.LockRes:
          setRoomLocked(true);
          break
        case ResponseType.UnlockRes:
          setRoomLocked(false)
          break
        default:
          console.log("Unhandled response type: " + res.type)
      }
    }

    setClient(ws);
  }

  const disconnect = () => {
    if (!client) {
      throw new Error("not connected");
    }

    client.close();
    console.log("closing");
  }

  const connect = () => {
    if (client) {
      throw new Error("already made client");
    }

    setResponses([]);

    console.log(responses.length);

    initClient();
  }

  const lockRoom = () => {
    if (!client) {
      throw new Error("not connected");
    }

    sendRequest(client, {
      type: RequestType.LockReq,
      send: Date.now(),
    } as WSRequest);
  }

  const unlockRoom = () => {
    if (!client) {
      throw new Error("not connected");
    }

    sendRequest(client, {
      type: RequestType.UnlockReq,
      send: Date.now(),
    } as WSRequest);
  }

  const sendGameEvent = () => {
    if (!client) {
      throw new Error("not connected");
    }

    sendRequest(client, {
      type: RequestType.GameEvtReq,
      id: uuidv4(),
      send: Date.now(),
      name: "my-game-event",
      payload: btoa(geInputRef.current!.value),
    } as WSRequest)

  }


  return (
    <div className="App">
      <br />
      <h1>Ata Net Room Tester</h1>
      <div className='flex-column'>
        <div className='card'>
          <h3>Connect to Server</h3>
          <div className='form-input'><input disabled={connected()} ref={serverRef} type="text" name="server" placeholder='WS Server'></input></div>
          <div className='form-input'><input disabled={connected()} ref={codeRef} type="text" name="room_code" placeholder='Room Code'></input></div>
          <div className='form-input'><input disabled={connected()} ref={nameRef} type="text" name="user_name" placeholder='Player Name'></input></div>

          <div className='form-input'>
            <label>Allow Create
              <input ref={createRef} type="checkbox" name="allow_create" disabled={connected()}></input>
            </label>
          </div>
          <div className='button-holder'>
            <button onClick={connect} disabled={connected()}>Connect</button>
            <button onClick={disconnect} disabled={!connected()}>Disconnect</button>
            <button onClick={lockRoom} disabled={!connected() || roomLocked}>Lock</button>
            <button onClick={unlockRoom} disabled={!connected() || !roomLocked}>Unlock</button>

          </div>
        </div>
      </div>
      <div className='flex-column'>
        <div className='card'>
          <h3>Send Game Event</h3>
          <div>
            <textarea ref={geInputRef} name="ge-payload" disabled={!connected()} />
          </div>
          <div>
            <button onClick={sendGameEvent} disabled={!connected()}>Send</button>
          </div>
        </div>
      </div>

      <div className='flex-column'>
        <div className='card'>
          <h3>Responses</h3>
          <button onClick={() => setResponses([])}>Clear</button>
          {renderResponses(responses)}
        </div>
      </div>
    </div>
  );
}

export default App;
