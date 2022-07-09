import axios from "axios";

import { API_ADDR } from "./consts";

export const getDef = async (sourceId, input, params) =>
  apiGet([sourceId, input], params);

export const getDocNames = async () => apiGet("docs");

export const getDoc = async (name) => apiGet(["docs", name]);

export const getSetups = async () => apiGet("setups");

const apiGet = async (pathElements, queryParams) => {
  const pathString = Array.isArray(pathElements)
    ? pathElements.join("/")
    : pathElements;

  let queryString = "";
  if (queryParams) {
    queryString =
      "?" +
      Object.keys(queryParams)
        .map((k) => `${k}=${queryParams[k]}`)
        .join("&");
  }

  const res = await axios.get(`${API_ADDR}/${pathString}${queryString}`);
  return res.data;
};
