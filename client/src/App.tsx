import { useState, useEffect, useRef } from 'react'
import './App.css'

function App() {
  const [websocket, setWebsocket] = useState<WebSocket>()
  const [socketStatus, setSocketStatus] = useState<string>('')
  const [message, setMessage] = useState<string>('')
  const [messages, setMessages] = useState<string[]>([])
  const countRef = useRef(0)
  const [second, setSecond] = useState(0)
  const [selectFPS, setSelectFPS] = useState(1)
  const [receive, setReceive] = useState(0)
  const [fps, setFps] = useState(0)
  const [rfps, setRFps] = useState(0)

  const [uid, setUid] = useState('')

  const socketServerUrl = import.meta.env.VITE_WS_SERVER_URL

  const onMessage = (event: MessageEvent<string>) => {
    setReceive((prev) => prev + 1)
    const broadcastData = JSON.parse(event.data)
    if (broadcastData.type === 'init') {
      setUid(broadcastData.body.uid)
      console.log(uid)
      console.log(broadcastData)
    }
    // console.log(event.data)
    setMessage(event.data)
    setMessages((prev) => [...prev, event.data])
  }

  const connect = (): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
      const socket = new WebSocket(socketServerUrl)
      socket.onopen = () => {
        console.log('connected')
        setSocketStatus('connected')
        setWebsocket(socket)
        setSecond(0)
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
      setSecond(0)
    })

    // useEffectのクリーンアップの中で、WebSocketのクローズ処理を実行
    return () => {
      if (websocket != null) {
        websocket.close()
        websocket.removeEventListener('message', onMessage)
      }
    }
  }, [])

  useEffect(() => {
    const intervalId = setInterval(() => {
      if (uid && websocket) {
        websocket.send(
          JSON.stringify({
            type: 'pos',
            body: {
              uid: uid.toString(),
              name: uid.toString(),
              x: countRef.current.toString(),
              y: countRef.current.toString(),
            },
          })
        )

        countRef.current = countRef.current + 1
      }
    }, Math.floor((1 / selectFPS) * 1000))
    return () => clearInterval(intervalId)
  }, [websocket, selectFPS])

  useEffect(() => {
    const intervalId = setInterval(() => {
      if (uid && websocket) {
        setSecond((prev) => prev + 1)
        // console.log(count)
      }
    }, 1000)
    return () => clearInterval(intervalId)
  }, [websocket])

  useEffect(() => {
    setFps(countRef.current / second)
    setRFps(receive)
    setReceive(0)
  }, [second])

  const handleChangeFPS = (e: React.ChangeEvent<HTMLSelectElement>) => {
    countRef.current = 0
    setSecond(0)
    setMessage('')
    setReceive(0)
    setSelectFPS(Number(e.target.value))
  }

  return (
    <>
      <h2 className="mb-10">接続状況：{socketStatus}</h2>
      <h2>最新メッセージ</h2>
      <div className="h-40 overflow-y-scroll bg-gray-100 px-5 py-2 rounded-xl my-5">
        {message}
      </div>
      <h2>メッセージ履歴</h2>
      <div className="h-40 overflow-y-scroll bg-gray-100 px-5 py-2 rounded-xl my-5">
        {messages}
      </div>
      <div className="bg-gray-100 px-5 py-2 rounded-xl my-5">
        <div>{fps} fps</div>
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
              countRef.current = 0
              setSecond(0)
              setMessage('')
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
              connect().then((socket) => {
                setWebsocket(socket)
                setSecond(0)
                countRef.current = 0
              })
            }}
          >
            再接続
          </button>
        </div>
      </div>
      <div className="bg-gray-100 px-5 py-2 rounded-xl my-5">
        {/* <div>{count} 回</div> */}
        <div>{rfps} 回</div>
      </div>
    </>
  )
}

export default App
