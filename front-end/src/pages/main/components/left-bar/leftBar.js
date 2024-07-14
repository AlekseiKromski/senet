import ChatList from "./chat-list/chatList";
import CurrentUser from "./current-user/currentUser";
import LeftBarStyle from "./leftBar.module.css"

export default function LeftBar() {
    return (
        <div className={LeftBarStyle.LeftBarWindow}>
            <CurrentUser/>
            <ChatList/>
        </div>
    )
}