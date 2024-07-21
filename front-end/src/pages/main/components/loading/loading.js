import LoadingStyle from "./loading.module.css";
import {Spin} from "antd";
import {LoadingOutlined} from "@ant-design/icons";

export default function Loading(){
    return (
        <div className={LoadingStyle.LoadingWindow}>
            <Spin className={LoadingStyle.Spin} indicator={<LoadingOutlined spin/>} size="large"/>
        </div>
    )
}