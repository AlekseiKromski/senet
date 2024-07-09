import {useDispatch, useSelector} from "react-redux";
import {useState} from "react";
import {useNavigate} from "react-router-dom";
import Cookies from 'js-cookie';
import LoginStyle from "./login.module.css"
import {Button, Checkbox, Form, Input, message} from "antd";
import {setUserId} from "../../store/application/application";

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
                dispatch(setUserId(res.data.uid))
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
        <div>
            {contextHolder}
            <Form
                name="basic"
                labelCol={{
                    span: 8,
                }}
                wrapperCol={{
                    span: 16,
                }}
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

                <Form.Item
                    wrapperCol={{
                        offset: 8,
                        span: 16,
                    }}
                >
                    <Button type="primary" htmlType="submit" loading={loader}>
                        Login
                    </Button>
                </Form.Item>
            </Form>
        </div>
    )
}