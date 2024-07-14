import {Outlet} from "react-router-dom"
import LayoutStyle from "./layout.module.css"

function Layout() {
    return (
        <div className={LayoutStyle.Layout}>
            <Outlet/>
        </div>
    );
}

export default Layout;
