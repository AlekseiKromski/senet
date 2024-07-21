import {createSlice} from "@reduxjs/toolkit";

const loadingSlice = createSlice({
        name: 'loading',
        initialState: {
            chatLoading: false,
        },
        reducers: {
            setChatLoading: (state, payload) => {
                let p = payload.payload
                state.chatLoading = p
            },
        },
    }
);
export const {setChatLoading} = loadingSlice.actions;

// this is for configureStore
export default loadingSlice.reducer;