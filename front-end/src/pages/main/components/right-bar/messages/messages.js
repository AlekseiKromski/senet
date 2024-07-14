import messagesStyle from "./messages.module.css";
import {useSelector} from "react-redux";

export default function Messages(){
    // REDUX store
    const messages = useSelector((state) => state.messages);
    const chat = useSelector((state) => state.chat);
    const user = useSelector((state) => state.user);

    return (
        <div className={messagesStyle.DialogMessagesWrapper}>
            {
                messages.messages && messages.messages[chat.currentChat.id].map(m => {
                    let isCurrentUser = chat.currentChat.users.find(u => user.user.id === m.sender_id)
                    return (
                        <div className={[messagesStyle.DialogMessage, isCurrentUser ? messagesStyle.DialogMessageCurrentUser : messagesStyle.DialogMessageOtherUser ].join(" ")}>
                            <b>{chat.currentChat && chat.currentChat.users.find(u => u.id === m.sender_id).username}</b>
                            <span>{m.message}</span>
                            <small>{m.created_at}</small>
                        </div>
                    )
                })
            }
        </div>
    )
}