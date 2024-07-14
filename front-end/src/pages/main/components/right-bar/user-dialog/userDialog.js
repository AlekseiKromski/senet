import {useSelector} from "react-redux";

export default function UserDialog(){
    // REDUX store
    const user = useSelector((state) => state.user);
    const chat = useSelector((state) => state.chat);

    return (
        <div>
            <b>{user.user.id ? chat.currentChat.users.find(u => u.id !== user.user.id).username : chat.currentChat.name}</b>
        </div>
    )
}