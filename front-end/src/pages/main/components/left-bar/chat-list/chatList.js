import {useDispatch, useSelector} from "react-redux";
import {setChat} from "../../../../../store/chat/chat";
import {setMessages} from "../../../../../store/messages/messages";
import ChatListStyle from "./chatlist.module.css";
import Moment from "react-moment";
import {useEffect} from "react";
import {setLastMessage} from "../../../../../store/chats/chats";
import {setChatLoading} from "../../../../../store/loading/loading";

// Represent list of chats
export default function ChatList() {
    const security_levels = Object.freeze({
        SERVER_PRIVATE_KEY: "SERVER_PRIVATE_KEY",
    });

    // REDUX store
    const dispatch = useDispatch()
    const chats = useSelector((state) => state.chats);
    const chat = useSelector((state) => state.chat);
    const application = useSelector((state) => state.application);
    const messages = useSelector((state) => state.messages);
    const user = useSelector((state) => state.user);

    const fetchMessages = async (cid) => {
        dispatch(setChatLoading(true))

        if (chat.currentChat && chat.currentChat.id === cid) {
            dispatch(setChat(null))
            setTimeout(() => {
                dispatch(
                    dispatch(setChatLoading(false))
                )
            }, 500)
            return
        }

        if (!messages.messages[cid] || messages.messages[cid].length < 14) {
            let response = await application.axios.get(`/api/chat/messages/get/${cid}/0`)
            dispatch(setMessages({
                cid: cid,
                messages: response.data,
                reverse: true
            }))
        }
        dispatch(setChat(
            chats.chats.find(chat => chat.id === cid)
        ))

        setTimeout(() => {
            dispatch(
                dispatch(setChatLoading(false))
            )
        }, 1000)
    }

    useEffect(() => {
        for (const [key, value] of Object.entries(messages.messages)) {
            if (value.length === 0) {
                continue
            }
            dispatch(setLastMessage({
                chatId: key,
                message: value[value.length - 1]
            }))
        }
    }, [messages]);

    return (
        <ul className={ChatListStyle.UsersList}>
            {
                chats.chats && chats.chats.map(c => {
                    if (c.security_level === security_levels.SERVER_PRIVATE_KEY) {
                        let userFromChat = c.users.find(u => u.id !== user.user.id)
                        return (
                            <li
                                onClick={() => fetchMessages(c.id)}
                                className={chat.currentChat && chat.currentChat.id === c.id ? ChatListStyle.Selected : ""}
                            >
                                <div className={ChatListStyle.UserChatBox}>
                                    <div className={ChatListStyle.ChatProfile}>
                                        <div className={ChatListStyle.UserProfileImage} style={
                                            {
                                                backgroundImage: `url('${process.env.REACT_APP_AXIOS_BASE_URL}/storage/images/${userFromChat.image}')`
                                            }
                                        }></div>
                                        {userFromChat.username}
                                    </div>
                                    {
                                        c.last_message && <div
                                            className={[ChatListStyle.LastMessage, ChatListStyle.LastMessageText].join(" ")}>
                                            <small>{c.last_message && c.last_message.message.substring(0, 10)}</small>
                                            <small>
                                                <Moment format="HH:mm">
                                                    {c.created_at}
                                                </Moment>
                                            </small>
                                        </div>
                                    }
                                </div>
                            </li>
                        )
                    }
                })
            }
        </ul>
    )
}