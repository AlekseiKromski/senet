import {createSlice} from '@reduxjs/toolkit';
import {axiosInstance} from "../axios";

const applicationSlice = createSlice({
        name: 'application',
        initialState: {
            axios: axiosInstance,
            websocket: null,
        },
        reducers: {
            initWs: (state, funcs) => {
                let ws = new WebSocket(
                    process.env.REACT_APP_WEBSOCKET_URL + "/ws/connect",
                )

                ws.onopen = e => funcs.payload.success(e)

                ws.onerror = e => funcs.payload.error(e)

                ws.onclose = e => funcs.payload.error(e)

                ws.onmessage = (e) => funcs.payload.events(e)

                state.websocket = ws
            }
        }
    }
);

// this is for dispatch
export const {initWs} = applicationSlice.actions;

// this is for configureStore
export default applicationSlice.reducer;