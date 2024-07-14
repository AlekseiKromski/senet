import CurrentUserStyle from "./curretUser.module.css"
import {useDispatch, useSelector} from "react-redux";
export default function CurrentUser() {
    // REDUX store
    const user = useSelector((state) => state.user);
    return  (
        <div className={CurrentUserStyle.CurrentUser}>
            <div className={CurrentUserStyle.CurrentUserProfileImage} style={
                {
                    backgroundImage: `url('${process.env.REACT_APP_AXIOS_BASE_URL}/storage/images/${user.user.image}')`
                }
            }></div>
            <div className={CurrentUserStyle.CurrentUserProfileBlock}>
                <b>[{user.user.username}] <small>{user.user.first_name} {user.user.second_name}</small></b>
                <small>
                    <div className={CurrentUserStyle.CurrentUserOnlineStatus}></div>
                     online
                </small>
            </div>
        </div>
    )
}