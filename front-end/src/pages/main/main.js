import {useEffect} from "react";
import {useDispatch, useSelector} from "react-redux";
import {message} from "antd";
import {initWs} from "../../store/application/application";
import MainStyle from "./main.module.css"
import LeftBar from "./components/left-bar/leftBar";
import Cookies from "js-cookie";
import {setUser} from "../../store/user/user";
import {useNavigate} from "react-router-dom";
import {setChats} from "../../store/chats/chats";
import RightBar from "./components/right-bar/rightBar";
import {setMessage} from "../../store/messages/messages";
import Loading from "./components/loading/loading";
import {setIsTyping} from "../../store/typing/typing";
import Default from "./components/default/default";

export default function Main() {
    const navigate = useNavigate()
    const dispatch = useDispatch();
    const application = useSelector((state) => state.application);
    const loading = useSelector((state) => state.loading);
    const user = useSelector((state) => state.user);
    const chat = useSelector((state) => state.chat);
    const [messageApi, contextHolder] = message.useMessage();

    useEffect(() => {
        // Get user from cookie (server have to always set user in cookie)
        let userJson = Cookies.get("user")
        if (!userJson) {
            messageApi.warning("Not authorized")
            return navigate("/login");
        }
        let user = JSON.parse(userJson)
        dispatch(setUser(user))

        let timeout = null

        dispatch(initWs(
            {
                success: () => {messageApi.success("Connected")},
                error: () => {messageApi.error("failed to connect")},
                events: (event) => {
                    let received = JSON.parse(event.data)
                    if (received.action === "STATE_MESSAGE_OK" || received.action === "INCOMING_MESSAGE") {
                        let payload = JSON.parse(received.payload)
                        dispatch(setMessage({
                            cid: payload.chat_id,
                            message: payload
                        }))
                    }

                    if (received.action === "TYPING") {
                        if (timeout){
                            clearTimeout(timeout)
                        }

                        timeout = setTimeout(() => {
                            dispatch(setIsTyping({
                                cid: cid,
                                state: false
                            }))
                        }, 1000)

                        let cid = received.payload
                        dispatch(setIsTyping({
                            cid: cid,
                            state: true
                        }))
                    }
                }
            }
        ))

        application.axios.get('/api/chat/get').then(
            response => {
                dispatch(setChats(response.data))
            }
        )
    }, []);

    return (
        <div className={MainStyle.MainWrapper}>
            {contextHolder}
            {
                user.user && <LeftBar/>
            }
            {
                !loading.chatLoading && !chat.currentChat && <Default/>
            }
            {
                loading.chatLoading && <Loading/>
            }
            {
                !loading.chatLoading && chat.currentChat && <RightBar/>
            }
        </div>
    )
}