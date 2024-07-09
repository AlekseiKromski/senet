import {Outlet} from "react-router-dom"
import LayoutStyle from "./layout.module.css"

function Layout() {
    return (
        <div>
            <div>
                <h1>Welcome to Senet</h1>
                <Outlet/>
            </div>
        </div>
    );
}

export default Layout;
