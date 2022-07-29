import React, { useRef } from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  const serverRef = useRef<HTMLInputElement>(null);
  const codeRef = useRef<HTMLInputElement>(null);
  const createRef = useRef<HTMLInputElement>(null);

  const connect = () => {
    const ws = new WebSocket(serverRef!.current!.value, "atanet_v1");
    const req = JSON.stringify({
      room_code: codeRef!.current!.value,
      create: createRef!.current!.checked,
      player_name: "Test Name",
      game_id: "React Tester",
      room_size: 16
    });

    
    const encoder = new TextEncoder();
    const decoder = new TextDecoder();

    ws.onopen = () => {
      console.log("send it!");
      const bytes = encoder.encode(req);
      ws.send(bytes)
    }
    ws.onmessage = (msg) => {
      console.log(msg);
      const res = JSON.parse(msg.data);
      console.dir(res);

      if(res.payload) {
        const payload = JSON.parse(atob(res.payload));
        console.dir(payload);
      }
    }

    // type WSJoinRequest struct {
    //   RoomCode    string `json:"room_code"`
    //   AllowCreate bool   `json:"create"`
    //   PlayerName  string `json:"player_name"`
    //   GameID      string `json:"game_id"`
    //   RoomSize    int    `json:"room_size"`
    // }
  }


  return (
    <div className="App">
      <br/>
      <h1>Ata Net Room Tester</h1>
      <div className='form-input'><input ref={serverRef} type="text" name="server" placeholder='WS Server'></input></div>
      <div className='form-input'><input ref={codeRef} type="text" name="room_code" placeholder='Room Code'></input></div>
      
      <div className='form-input'>
        <label>Allow Create
          <input ref={createRef} type="checkbox" name="allow_create"></input>
        </label>
      </div>
      <button onClick={connect}>Send</button>
      {/* <div className='form-input'><input type="text" name="server" placeholder='WS Server'></input></div> */}
      {/* <div className='form-input'><input type="text" name="server" placeholder='WS Server'></input></div> */}
    </div>
  );
}

export default App;
