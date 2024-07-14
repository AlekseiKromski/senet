import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import reportWebVitals from './reportWebVitals';
import LayoutLogin from "./layout/login/layout"
import LayoutMain from "./layout/main/layout"
import {Provider} from "react-redux";
import {createHashRouter, RouterProvider} from "react-router-dom";
import Login from "./pages/login/login";
import store from "./store/store"
import Main from "./pages/main/main";

const router = createHashRouter([
    {
        element: <LayoutMain/>,
        children: [
            {
                path: "/", element: <Main/>,
            },
        ],
        path: "/"
    },
    {
        element: <LayoutLogin/>,
        children: [
            {
                path: "/login", element: <Login/>,
            },
        ],
    }
])

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <Provider store={store}>
        <RouterProvider router={router}/>
    </Provider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
