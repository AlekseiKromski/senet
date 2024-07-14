import MainStyle from "../../../main.module.css";
import {useDispatch, useSelector} from "react-redux";
import {setChat} from "../../../../../store/chat/chat";
import {setMessages} from "../../../../../store/messages/messages";
import ChatListStyle from "./chatlist.module.css";

// Represent list of chats
export default function ChatList() {
    const security_levels = Object.freeze({
        SERVER_PRIVATE_KEY: "SERVER_PRIVATE_KEY",
    });

    // REDUX store
    const dispatch = useDispatch()
    const chats = useSelector((state) => state.chats);
    const application = useSelector((state) => state.application);
    const messages = useSelector((state) => state.messages);
    const user = useSelector((state) => state.user);

    const fetchMessages = async (cid) => {
        if (!messages.messages[cid]) {
            let response = await application.axios.get('/api/chat/messages/get/' + cid)
            dispatch(setMessages({
                cid: cid,
                messages: response.data
            }))
        }
        dispatch(setChat(
            chats.chats.find(chat => chat.id === cid)
        ))
    }

    return (
        <ul className={ChatListStyle.UsersList}>
            {
                chats.chats && chats.chats.map(c => {
                    if (c.security_level === security_levels.SERVER_PRIVATE_KEY) {
                        let userFromChat = c.users.find(u => u.id !== user.user.id)
                        return (
                            <li
                                onClick={() => fetchMessages(c.id)}
                            >
                                <div className={ChatListStyle.UserProfileImage} style={
                                    {
                                        backgroundImage: `url('${process.env.REACT_APP_AXIOS_BASE_URL}/storage/images/${userFromChat.image}')`
                                    }
                                }></div>
                                {userFromChat.username}
                            </li>
                        )
                    }
                })
            }
        </ul>
    )
}