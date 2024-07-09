import {createSlice} from '@reduxjs/toolkit';
import axios from "axios"

const applicationSlice = createSlice({
        name: 'axios',
        initialState: {
            axios: function () {
                let instance = axios.create({
                    baseURL: process.env.REACT_APP_AXIOS_BASE_URL,
                    timeout: 1000,
                    headers: {},
                    withCredentials: true
                })
                instance.defaults.timeout = 15000;

                instance.interceptors.response.use(
                    (response) => {
                        return response;
                    },
                    (error) => {
                        if (error.response && error.response.status === 401) {
                            // Use router.push() to navigate to the login screen
                            document.location = "/#login"
                            // Throw an exception to stop further execution
                            return Promise.reject(error.response.data.message);
                        }
                        // Handle other errors here
                        return Promise.reject(error);
                    }
                );

                return instance
            } (),
            userId: null
        },
        reducers: {
            setUserId: (state, id) => {
                state.userId = id.payload
            }
        }
    }
);

// this is for dispatch
export const {setUserId} = applicationSlice.actions;

// this is for configureStore
export default applicationSlice.reducer;