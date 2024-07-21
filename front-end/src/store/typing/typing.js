import {createSlice} from "@reduxjs/toolkit";

const typingSlice = createSlice({
        name: 'typing',
        initialState: {
            isTyping: {}
        },
        reducers: {
            setIsTyping: (state, payload) => {
                let p = payload.payload
                state.isTyping[p.cid] = p.state
            }
        },
    }
);
export const {setIsTyping} = typingSlice.actions;

// this is for configureStore
export default typingSlice.reducer;