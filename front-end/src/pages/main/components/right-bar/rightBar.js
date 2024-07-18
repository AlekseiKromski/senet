import RightBarStyle from "./rightBar.module.css";
import UserDialog from "./user-dialog/userDialog";
import Messages from "./messages/messages";
import SendText from "./send-text/sendText";
import {useState} from "react";

export default function RightBar(){

    const [send, setSend] = useState(false)

    return (
        <div className={RightBarStyle.DialogWindow}>
            <UserDialog/>
            <Messages send={send}/>
            <SendText setSend={setSend} send={send}/>
        </div>
    )
}