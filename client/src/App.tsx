import { useState, useRef, useEffect, ChangeEvent } from 'react'
import './App.css'
// import * as WebSocket from 'websocket'

function App() {
  const [websocket, setWebsocket] = useState<WebSocket>()
  const [socketStatus, setSocketStatus] = useState<string>('')
  const [message, setMessage] = useState<string[]>([])
  const socketRef = useRef<WebSocket>()
  const [count, setCount] = useState(0)
  const [second, setSecond] = useState(0)
  const [selectFPS, setSelectFPS] = useState(1)
  const [uid, setUid] = useState('')

  const socketServerUrl = 'ws://localhost:8080/ws'

  const onMessage = (event: MessageEvent<string>) => {
    const broadcastData = JSON.parse(event.data)
    if (broadcastData.type === 'init') {
      setUid(broadcastData.body.uid)
      console.log(uid)
      console.log(broadcastData)
    }
    console.log(event.data)
    setMessage((prev) => [...prev, event.data])
  }

  const connect = (): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
      const socket = new WebSocket(socketServerUrl)
      socket.onopen = () => {
        console.log('connected')
        setSocketStatus('connected')
        setWebsocket(socket)
        resolve(socket)
      }
      socket.onmessage = onMessage
      socket.onclose = () => {
        console.log('reconnecting...')
        setSocketStatus('reconnecting...')
        connect()
      }
      socket.onerror = (err: Event) => {
        console.log('connection error:', err)
        reject(err)
      }
    })
  }

  useEffect(() => {
    connect().then((socket) => {
      setWebsocket(socket)
    })

    // #3.useEffectのクリーンアップの中で、WebSocketのクローズ処理を実行
    return () => {
      if (websocket != null) {
        websocket.close()
        websocket.removeEventListener('message', onMessage)
      }
    }
  }, [])

  useEffect(() => {
    const intervalId = setInterval(() => {
      if (uid && websocket && websocket?.readyState === websocket?.OPEN) {
        websocket.send(
          JSON.stringify({
            type: 'pos',
            body: {
              uid: uid.toString(),
              name: uid.toString(),
              x: count.toString(),
              y: count.toString(),
            },
          })
        )
      }
      setCount((prev) => prev + 1)
    }, Math.floor((1 / selectFPS) * 1000))
    return () => clearInterval(intervalId)
  }, [count])

  useEffect(() => {
    const intervalId = setInterval(() => {
      setSecond((prev) => prev + 1)
    }, 1000)
    return () => clearInterval(intervalId)
  }, [second])

  const handleChangeFPS = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setCount(0)
    setSecond(0)
    setMessage([])
    setSelectFPS(Number(e.target.value))
  }

  return (
    <>
      <div className="text-xl">{socketStatus}</div>
      <div className="h-40 overflow-y-scroll bg-gray-100 px-5 py-2 rounded-xl my-5">
        {message}
      </div>
      <div className="bg-gray-100 px-5 py-2 rounded-xl my-5">
        <div>{count} 回</div>
        <div>{Math.floor(count / second)} fps</div>
      </div>
      <div className="bg-gray-100 px-5 py-2 rounded-xl my-5">
        <div className="my-3 flex w-fit mx-auto">
          <p>select fps: </p>
          <select
            className="bg-inherit rounded-md mx-2"
            onChange={(e) => handleChangeFPS(e)}
          >
            <option value={1}>1</option>
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={30}>30</option>
          </select>
        </div>
        <div className="my-3">
          <button
            className="rounded-md bg-sky-500 text-white px-5 py-1"
            type="button"
            onClick={() => {
              setCount(0)
              setSecond(0)
              setMessage([])
            }}
          >
            リセット
          </button>
        </div>
        <div className="my-3">
          <button
            className="rounded-md bg-teal-500 text-white px-5 py-1"
            type="button"
            onClick={() => {
              websocket?.close()
            }}
          >
            再接続
          </button>
        </div>
      </div>
    </>
  )
}

export default App
