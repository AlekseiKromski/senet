import {useState} from "react";
import {Avatar, Button, Input, List, Modal, Skeleton} from "antd";
import AddChatUserStyle from "./addChatUser.module.css"
import {useDispatch, useSelector} from "react-redux";
import {setChats} from "../../../../../store/chats/chats";
export default function AddChatWithUser(){

    // REDUX
    const dispatch = useDispatch()
    const application = useSelector((state) => state.application);
    const user = useSelector((state) => state.user);
    const chats = useSelector((state) => state.chats);

    // STATE
    const [open, setOpen] = useState(false);
    const [confirmLoading, setConfirmLoading] = useState(false);
    const [users, setUsers] = useState([])
    const [initLoading, setInitLoading] = useState(true);
    const [initLoadingCreateChat, setInitLoadingCreateChat] = useState(false)

    const showModal = () => {
        setOpen(true);
    };
    const handleOk = () => {
        setConfirmLoading(true);
        setTimeout(() => {
            setOpen(false);
            setConfirmLoading(false);
        }, 500);
    };

    const handleCancel = (e) => {
        setOpen(false)
        setUsers([])
    };

    const searchUsersByUsername = (e) => {

        let username = e.target.value
        if (username.length == 0) {
            return
        }
        setInitLoading(true)

        application.axios.get(`/api/users/${username}`).then(
            response => {
                setTimeout(() => {
                    setUsers([...response.data])
                    setInitLoading(false)
                }, 1000)
            }
        ).catch(e => {
            console.error(e)
        })
    }

    const createChatWithUser = (id) => {
        setInitLoadingCreateChat(true)
        application.axios.post(`/api/chat/create`, {
            "user1":  user.user.id,
            "user2": id,
            "chat_type":  "PRIVATE",
            "security_level": "SERVER_PRIVATE_KEY"
        }).then(
            response => {
                dispatch(setChats([...chats.chats, response.data]))
                setTimeout(() => {
                    setInitLoadingCreateChat(false)
                }, 1000)
            }
        )
    }
    return (
        <div className={AddChatUserStyle.AddChatUserBlock}>
            <Button type="primary" onClick={showModal}>
                Find user
            </Button>
            <Modal
                title="Find user"
                open={open}
                onOk={handleOk}
                confirmLoading={confirmLoading}
                onCancel={handleCancel}
            >
                <Input
                    onChange={searchUsersByUsername}
                    placeholder="Write some username..."
                />
                {
                    users.length !== 0 && <List
                        loading={initLoading}
                        itemLayout="horizontal"
                        dataSource={users}
                        renderItem={(item) => (
                            <List.Item
                                actions={[
                                    chats.chats.find(c => {
                                        return c.users.find(u => u.id === item.id)
                                    }) ? <b>Exists</b> : <Button onClick={() => {createChatWithUser(item.id)}} color="primary" loading={initLoadingCreateChat} >+</Button>
                                ]}
                            >
                                <Skeleton avatar title={false} loading={item.loading} active>
                                    <List.Item.Meta
                                        avatar={<Avatar src={`${process.env.REACT_APP_AXIOS_BASE_URL}/storage/images/${item.image}`} />}
                                        title={item.username}
                                    />
                                </Skeleton>
                            </List.Item>
                        )}
                    />
                }
            </Modal>
        </div>
    );
}