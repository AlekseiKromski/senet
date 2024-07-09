import { configureStore } from '@reduxjs/toolkit';
import applicationReducer from "./application/application";

export default configureStore({
    reducer: {
        application: applicationReducer
    },
});