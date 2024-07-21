import messagesStyle from "./messages.module.css";
import {useDispatch, useSelector} from "react-redux";
import Scrollbars from "react-custom-scrollbars-2";
import {useEffect, useRef, useState} from "react";
import {Button, Spin} from "antd";
import Moment from "react-moment";
import {setMessages} from "../../../../../store/messages/messages";
import {LoadingOutlined} from "@ant-design/icons";

export default function Messages({send}) {
    // REDUX store
    const dispatch = useDispatch()
    const messages = useSelector((state) => state.messages);
    const chat = useSelector((state) => state.chat);
    const user = useSelector((state) => state.user);
    const application = useSelector((state) => state.application);

    // REF
    const scrollbarRef = useRef(null);

    // STATE
    const [back, setBack] = useState(false)
    const [loading, setLoading] = useState(false)
    const [offset, setOffset] = useState(15)
    const [shouldLoad, setShouldLoad] = useState(true)

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
        setShouldLoad(true)
    }, [chat.currentChat]);

    useEffect(() => {
        // When not enough messages, no need to use back
        if (scrollbarRef.current.getScrollHeight() < 1000) {
            scrollbarRef.current.scrollToBottom();
            return;
        }

        let ms = messages.messages[chat.currentChat.id]
        if (!ms || ms.length === 0){
            return
        }
        let percentageFromBottom = getPercentage()
        if (percentageFromBottom.toFixed(2) < 40.0) {
            scrollbarRef.current.scrollToBottom();
            setBack(false)
        } else {
            if (user.user.id === ms[ms.length - 1].sender_id ) {
                scrollbarRef.current.scrollToBottom();
            }
            setBack(true)
        }
    }, [messages.messages])

    // Functions
    const onScrollBack = () => {
        let percentageFromBottom = getPercentage()
        if (percentageFromBottom.toFixed(2) > 40.0) {
            setBack(true)
            return
        }

        setBack(false)
        return
    }

    const onScrollLoadMessages = () => {
        if (!shouldLoad) {
            return;
        }

        let percentageFromBottom = getPercentage()
        if (percentageFromBottom.toFixed(2) < 80.0 || loading) {
            return
        }

        setLoading(true)

        const cid = chat.currentChat.id

        if (!messages.messages[cid]) {
            setLoading(false)
            return
        }

        application.axios.get(`/api/chat/messages/get/${cid}/${offset}`)
            .then(
                response => {
                    if (response.data.length === 0) {
                        setShouldLoad(false)
                        setLoading(false)
                        return
                    }
                    const scrollTop = scrollbarRef.current.getScrollTop();
                    const scrollHeight = scrollbarRef.current.getScrollHeight();

                    dispatch(setMessages({
                        cid: cid,
                        messages: [...response.data.reverse(), ...messages.messages[cid]],
                        reverse: false,
                    }))

                    setLoading(false)
                    setOffset(offset + 15)
                    setTimeout(() => {
                        // Calculate the new scrollTop to maintain the same relative scroll position
                        const newScrollHeight = scrollbarRef.current.getScrollHeight();
                        const heightDifference = newScrollHeight - scrollHeight;
                        scrollbarRef.current.scrollTop(scrollTop + heightDifference);
                    }, 0);
                }
            )
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
                onScroll={() => {
                    onScrollBack()
                    onScrollLoadMessages()
                }}
                renderView={props => <div {...props} className={messagesStyle.ScrollbarView}/>}
                ref={scrollbarRef}
            >
                {loading && <Spin className={messagesStyle.Spin} indicator={<LoadingOutlined spin/>} size="large"/>}
                {
                    messages.messages && messages.messages[chat.currentChat.id].map(m => {
                        let isCurrentUser = chat.currentChat.users.find(u => user.user.id === m.sender_id)
                        return (
                            <div
                                className={[messagesStyle.DialogMessage, isCurrentUser ? messagesStyle.DialogMessageCurrentUser : messagesStyle.DialogMessageOtherUser].join(" ")}>
                                <b>{chat.currentChat && chat.currentChat.users.find(u => u.id === m.sender_id).username}</b>
                                {m.message}
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