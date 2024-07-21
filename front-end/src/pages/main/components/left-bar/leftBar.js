import ChatList from "./chat-list/chatList";
import CurrentUser from "./current-user/currentUser";
import LeftBarStyle from "./leftBar.module.css"
import AddChatWithUser from "./add-chat-with-user/addChatWithUser";

export default function LeftBar() {
    return (
        <div className={LeftBarStyle.LeftBarWindow}>
            <div>
                <CurrentUser/>
                <ChatList/>
            </div>
            <AddChatWithUser/>
        </div>
    )
}