import {createSlice} from '@reduxjs/toolkit';

const userSlice = createSlice({
        name: 'user',
        initialState: {
            user: null,
        },
        reducers: {
            setUser(state, user) {
                state.user = user.payload
            }
        }
    }
);

// this is for dispatch
export const {setUser} = userSlice.actions;

// this is for configureStore
export default userSlice.reducer;