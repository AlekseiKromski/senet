import SendTextStyle from "./sendText.module.css";
import TextArea from "antd/es/input/TextArea";
import {Button} from "antd";
import {useState} from "react";
import {useSelector} from "react-redux";

export default function SendText() {
    const [text, setText] = useState("")

    // REDUX STORE
    const chat = useSelector((state) => state.chat);
    const application = useSelector((state) => state.application);
    const user = useSelector((state) => state.user);

    const sendTextViaWs = () => {
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
    }

    return (
        <div className={SendTextStyle.DialogInput}>
            <TextArea onChange={(e) => setText(e.target.value)} value={text} placeholder="Autosize height based on content lines" autoSize/>
            <Button onClick={sendTextViaWs}>Send</Button>
        </div>
    )
}