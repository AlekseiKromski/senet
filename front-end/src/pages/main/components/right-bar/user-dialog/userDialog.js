import {useSelector} from "react-redux";
import UserDialogStyle from "./userDialog.module.css"

export default function UserDialog(){
    // REDUX store
    const user = useSelector((state) => state.user);
    const chat = useSelector((state) => state.chat);
    const typing = useSelector((state) => state.typing);

    return (
        <div className={UserDialogStyle.UserDialog}>
            <b>{user.user.id ? chat.currentChat.users.find(u => u.id !== user.user.id).username : chat.currentChat.name}</b>
            {
                typing.isTyping[chat.currentChat.id] ?
                <small className={UserDialogStyle.Status}>Typing...</small> : <small className={UserDialogStyle.Status}>Pending online status...</small>
            }
        </div>
    )
}