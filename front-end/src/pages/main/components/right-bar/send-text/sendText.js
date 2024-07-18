import SendTextStyle from "./sendText.module.css";
import TextArea from "antd/es/input/TextArea";
import {Button} from "antd";
import {useState} from "react";
import {useSelector} from "react-redux";

export default function SendText({setSend}) {
    const [text, setText] = useState("")

    // REDUX STORE
    const chat = useSelector((state) => state.chat);
    const application = useSelector((state) => state.application);
    const user = useSelector((state) => state.user);

    const sendTextViaWs = (e) => {
        if (text.length === 0) {
            return
        }
        application.websocket.send(
            JSON.stringify(
                {
                    action: "SENT_MESSAGE",
                    payload: JSON.stringify({
                        cid: chat.currentChat.id,
                        to: chat.currentChat.users.find(u => u.id !== user.user.id).id,
                        message: text
                    })
                }
            )
        )

        setText("")
        setSend(true)
        setTimeout(() => {
            setSend(false)
        }, 2000)
    }

    return (
        <div className={SendTextStyle.DialogInput}>
            <TextArea
                onChange={(e) => {
                    setText(e.target.value)
                }}
                value={text}
                autoSize
            />
            <Button onClick={sendTextViaWs}>Send</Button>
        </div>
    )
}