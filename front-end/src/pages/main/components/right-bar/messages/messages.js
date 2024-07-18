import messagesStyle from "./messages.module.css";
import {useSelector} from "react-redux";
import Scrollbars from "react-custom-scrollbars-2";
import {useEffect, useRef, useState} from "react";
import {Button} from "antd";
import Moment from "react-moment";

export default function Messages({send}) {
    // REDUX store
    const messages = useSelector((state) => state.messages);
    const chat = useSelector((state) => state.chat);
    const user = useSelector((state) => state.user);

    // REF
    const scrollbarRef = useRef(null);

    // State
    const [back, setBack] = useState(false)

    // Functions
    const getPercentage = () => {
        let scrollTop = scrollbarRef.current.getScrollTop();
        let scrollHeight = scrollbarRef.current.getScrollHeight();
        let clientHeight = scrollbarRef.current.getClientHeight();

        // Calculate the distance from the bottom
        let scrollFromBottom = scrollHeight - scrollTop - clientHeight;

        // Calculate the percentage from the bottom
        return (scrollFromBottom / (scrollHeight - clientHeight)) * 100;
    }

    useEffect(() => {
        scrollbarRef.current.scrollToBottom();
    }, []);

    useEffect(() => {
        if (send) {
            scrollbarRef.current.scrollToBottom();
        }
    }, [send])

    useEffect(() => {
        let percentageFromBottom = getPercentage()
        if (percentageFromBottom.toFixed(2) < 40.0) {
            scrollbarRef.current.scrollToBottom();
            setBack(false)
        } else {
            setBack(true)
        }
    }, [messages.messages])

    // Functions
    const onScroll = () => {
        let percentageFromBottom = getPercentage()
        if (percentageFromBottom.toFixed(2) > 40.0) {
            return
        }

        if (!back) {
            return
        }

        setBack(false)
    }

    return (
        <div className={messagesStyle.DialogMessagesWrapper}>
            {
                back && <Button onClick={() => {
                    scrollbarRef.current.scrollToBottom();
                    setBack(false)
                }} className={messagesStyle.BackButton} type="primary">
                    Back
                </Button>
            }
            <Scrollbars
                onScroll={onScroll}
                renderView={props => <div {...props} className={messagesStyle.ScrollbarView}/>}
                ref={scrollbarRef}
            >
                {
                    messages.messages && messages.messages[chat.currentChat.id].map(m => {
                        let isCurrentUser = chat.currentChat.users.find(u => user.user.id === m.sender_id)
                        return (
                            <div
                                className={[messagesStyle.DialogMessage, isCurrentUser ? messagesStyle.DialogMessageCurrentUser : messagesStyle.DialogMessageOtherUser].join(" ")}>
                                <b>{chat.currentChat && chat.currentChat.users.find(u => u.id === m.sender_id).username}</b>
                                <span>{m.message}</span>
                                <small>
                                    <Moment format="HH:mm">
                                        {m.created_at}
                                    </Moment>
                                </small>
                            </div>
                        )
                    })
                }
            </Scrollbars>
        </div>
    )
}