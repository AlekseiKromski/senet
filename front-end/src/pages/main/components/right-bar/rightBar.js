import RightBarStyle from "./rightBar.module.css";
import UserDialog from "./user-dialog/userDialog";
import Messages from "./messages/messages";
import SendText from "./send-text/sendText";

export default function RightBar(){
    return (
        <div className={RightBarStyle.DialogWindow}>
            <UserDialog/>
            <Messages/>
            <SendText/>
        </div>
    )
}