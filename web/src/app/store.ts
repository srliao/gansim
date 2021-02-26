import { configureStore, ThunkAction, Action } from "@reduxjs/toolkit";
import profileReducer from "features/profile/profileSlice";

export const store = configureStore({
  reducer: {
    profile: profileReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
