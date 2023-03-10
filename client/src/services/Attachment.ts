import axios from "axios";
import { baseURL, makeRequest } from "./makeRequest";

export const uploadAttachment = async (
  file: File,
  msgId: string,
  recipient: string,
  isRoomMsg: boolean
) => {
  const metaData = {
    ID: msgId,
    mime_type: file.type,
    name: file.name,
    size: file.size,
  };
  await axios({
    method: "POST",
    url: `${baseURL}/api/attachment/meta?${
      isRoomMsg ? "channel_id" : "uid"
    }=${recipient}`,
    data: metaData,
    withCredentials: true,
  });
  // Split attachment into 4mb chunks
  let fileUploadChunks: Promise<ArrayBuffer>[] = [];
  let startPointer = 0;
  let endPointer = file.size;
  while (startPointer < endPointer) {
    let newStartPointer = startPointer + 4 * 1024 * 1024;
    fileUploadChunks.push(
      new Blob([file.slice(startPointer, newStartPointer)]).arrayBuffer()
    );
    startPointer = newStartPointer;
  }
  // Upload chunks
  for await (const data of fileUploadChunks) {
    await makeRequest(
      `${baseURL}/api/attachment/chunk/${msgId}?${
        isRoomMsg ? "channel_id" : "uid"
      }=${recipient}`,
      {
        withCredentials: true,
        method: "POST",
        headers: { "Content-Type": "application/octet-stream" },
        data,
      }
    );
  }
};
