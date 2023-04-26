import { useState, useRef, useEffect, ChangeEvent } from 'react'
import './App.css'

function App() {
  const [message, setMessage] = useState<string[]>([])
  const socketRef = useRef<WebSocket>()
  const [count, setCount] = useState(0)
  const [second, setSecond] = useState(0)
  const [selectFPS, setSelectFPS] = useState(1)

  // #0.WebSocket関連の処理は副作用なので、useEffect内で実装
  useEffect(() => {
    // #1.WebSocketオブジェクトを生成しサーバとの接続を開始
    const websocket = new WebSocket('ws://localhost:8080/ws?test=8967124381964')
    socketRef.current = websocket

    // #2.メッセージ受信時のイベントハンドラを設定
    const onMessage = (event: MessageEvent<string>) => {
      setMessage(prev => [...prev, event.data])
    }
    websocket.addEventListener('message', onMessage)

    // #3.useEffectのクリーンアップの中で、WebSocketのクローズ処理を実行
    return () => {
      websocket.close()
      websocket.removeEventListener('message', onMessage)
    }
  }, [])

  useEffect(() => {
     const intervalId = setInterval(() => {
      socketRef.current?.send('送信メッセージ')
      setCount(prev => prev + 1)
    }, Math.floor(1/selectFPS*1000));
    return () => clearInterval(intervalId);
  }, [count])

  useEffect(() => {
    const intervalId = setInterval(() => {
      setSecond(prev => prev + 1)
   }, 1000);
   return () => clearInterval(intervalId);
 }, [second])

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setCount(0)
    setSecond(0)
    setSelectFPS(Number(e.target.value))
  }

  return (
    <>
      <div className='h-60 overflow-y-scroll'>
      </div>
      <div>{count} 回</div>
      <div>{Math.floor(count/second)} fps</div>
      <div>
        <select onChange={(e) => handleChange(e)}>
          <option value={0}>1</option>
          <option value={10}>10</option>
          <option value={20}>20</option>
          <option value={30}>30</option>
        </select>
      </div>
      <button
        type="button"
        onClick={() => {
          setCount(0)
          setSecond(0)
        }}
      >
        リセット
      </button>
    </>
  )
}

export default App
