import {useEffect, useRef, useState} from "react";
import {useSelector} from "react-redux";
import TextArea from "antd/es/input/TextArea";
import {Button, message} from "antd";

export default function Main() {

    const application = useSelector((state) => state.application);
    const [uid, setUid] = useState(null)
    const [chat, setChat] = useState(null)
    const [chats, setChats] = useState([])
    const [messages, setMessages] = useState([])
    const [text, setText] = useState("")
    const [messageApi, contextHolder] = message.useMessage();
    const ws = useRef(null);

    const loadMessages = (chatId) => {
        application.axios.get('/api/chat/messages/get/' + chatId).then(
            response => {
                setMessages(response.data.reverse())
                setChat(chats.find(c => c.id === chatId))
            }
        )
    }

    const sendTextViaWs = () => {
        ws.current.send(
            JSON.stringify(
                {
                    action: "SENT_MESSAGE",
                    payload: JSON.stringify({
                        cid: chat.id,
                        to: chat.users.find(u => u.id !== uid).id,
                        message: text
                    })
                }
            )
        )
    }

    const connect = () => {
        ws.current = new WebSocket(
            process.env.REACT_APP_WEBSOCKET_URL + "/ws/connect",
        )
        ws.current.onopen = e => {
            messageApi.open({
                type: 'success',
                content: `Connected`,
            });
        }

        ws.current.onerror = () => {
            messageApi.open({
                type: 'error',
                content: `Cannot connect`,
            });
        }

        ws.current.onclose = () => {
            messageApi.open({
                type: 'warning',
                content: `Connection closed`,
            });
        }
    }

    useEffect(() => {
        application.axios.get('/api/chat/get').then(
            response => {
                setChats(response.data)
            }
        )
        connect()
        return () => {
            ws.current.close()
        }
    }, []);

    useEffect(() => {
        if (!ws.current) return;

        ws.current.onmessage = (event) => {
            let received = JSON.parse(event.data)
            if (received.action === "USER_ID") {
                setUid(received.payload)
                return
            }
            if (received.action === "STATE_MESSAGE_OK" || received.action === "INCOMING_MESSAGE") {
                let payload = JSON.parse(received.payload)
                setMessages([...messages, payload])
                return
            }
        }
    }, [messages, uid])

    return (
        <div>
            {contextHolder}
            <ul>
                {
                    chats && chats.map(c => {
                        if (c.security_level === "SERVER_PRIVATE_KEY") {
                            return (
                                <li
                                    onClick={() => loadMessages(c.id)}
                                >{
                                    uid ? c.users.find(u => u.id !== uid).username : c.name
                                }</li>
                            )
                        }

                    })
                }
            </ul>
            <hr/>
            <div>
                {
                    chat && messages && messages.map(m => {
                        return (
                            <div className="">
                                <b>[{chat && chat.users.find(u => u.id === m.sender_id).username}]</b>: {m.message} <small>({m.created_at})</small>
                            </div>
                        )
                    })
                }
            </div>
            <hr/>
            {
                chat && <div>
                    <TextArea onChange={(e) => setText(e.target.value)} placeholder="Autosize height based on content lines"
                              autoSize/>
                    <Button onClick={sendTextViaWs}>Send</Button>
                </div>
            }
        </div>
    )
}