import {useDispatch, useSelector} from "react-redux";
import {useState} from "react";
import {useNavigate} from "react-router-dom";
import LoginStyle from "./login.module.css"
import {Button, Checkbox, Form, Input, message} from "antd";
import {setUser} from "../../store/user/user";
import Cookies from "js-cookie";

export default function Login() {

    let navigate = useNavigate();
    const dispatch = useDispatch();
    const application = useSelector((state) => state.application);
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [loader, setLoader] = useState(false)
    const [messageApi, contextHolder] = message.useMessage();

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
    };

    const login = () => {
        setLoader(true)

        application.axios.post(`/api/auth`, {
            username: username,
            password: password,
            type: "cookie"
        }, {
            withCredentials: true
        })
            .then(res => {
                messageApi.open({
                    type: 'success',
                    content: `Successful login`,
                });

                // Get user from cookie (server have to always set user in cookie)
                let userJson = Cookies.get("user")
                if (!userJson || userJson.length === 0) {
                    throw new Error('Cannot get user from cookies');
                }
                let user = JSON.parse(userJson)
                dispatch(setUser(user))

                setTimeout(() => {
                    setLoader(false)
                    navigate("/")
                }, 1000)
            })
            .catch(e => {
                messageApi.open({
                    type: 'error',
                    content: `Cannot login ${e}`,
                });
                setTimeout(() => setLoader(false), 1000)
            })
    }

    return (
        <div className={LoginStyle.FormWrapper}>
            {contextHolder}
            <div className={LoginStyle.Form}>
                <h1>Login</h1>
                <Form
                    name="basic"
                    layout="vertical"
                    style={{
                        maxWidth: 600,
                    }}
                    initialValues={{
                        remember: true,
                    }}
                    autoComplete="off"
                    onFinish={login}
                    onFinishFailed={onFinishFailed}
                >
                    <Form.Item
                        label="Username"
                        name="username"
                        rules={[
                            {
                                required: true,
                                message: 'Please input your username!',
                            },
                        ]}
                    >
                        <Input onChange={(e) => setUsername(e.target.value)} />
                    </Form.Item>

                    <Form.Item
                        label="Password"
                        name="password"
                        rules={[
                            {
                                required: true,
                                message: 'Please input your password!',
                            },
                        ]}
                    >
                        <Input.Password onChange={(e) => setPassword(e.target.value)} />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={loader}>
                            Login
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </div>
    )
}