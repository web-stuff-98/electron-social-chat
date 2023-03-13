import { reactive, ref } from "vue";

export const userMedia = ref<MediaStream>();
export const displayMedia = ref<MediaStream>();

export const userMediaProperties = reactive({
  audio: false,
  video: false,
});
export const displayMediaActive = ref(false);

function addUserMediaEventListeners() {
  userMedia.value?.addEventListener("addtrack", (t) => {
    if (t.track.kind === "audio") userMediaProperties.audio = true;
    if (t.track.kind === "video") userMediaProperties.video = true;
  });
  userMedia.value?.addEventListener("removetrack", (t) => {
    if (t.track.kind === "audio") userMediaProperties.audio = false;
    if (t.track.kind === "video") userMediaProperties.video = false;
  });
}

export async function openMic() {
  if (!userMedia.value) {
    userMedia.value = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: true,
    });
    userMediaProperties.video = true;
    userMediaProperties.audio = true;
    addUserMediaEventListeners();
  } else {
    const media = await navigator.mediaDevices.getUserMedia({
      video: userMediaProperties.video,
      audio: true,
    });
    userMediaProperties.audio = true;
    if (media) {
      media.getAudioTracks().forEach((track) => {
        userMedia.value?.addTrack(track);
        userMediaProperties.audio = true;
      });
    }
  }
}
export function muteMic() {
  if (userMedia.value) {
    userMedia.value.getAudioTracks().forEach((track) => {
      userMedia.value?.removeTrack(track);
      userMediaProperties.audio = false;
    });
  }
}

export async function openCamera() {
  if (!userMedia.value) {
    userMedia.value = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: userMediaProperties.audio,
    });
    userMediaProperties.video = true;
    addUserMediaEventListeners();
  } else {
    const media = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: userMediaProperties.audio,
    });
    if (media) {
      media.getVideoTracks().forEach((track) => {
        userMedia.value?.addTrack(track);
        userMediaProperties.video = true;
      });
    }
  }
}
export function closeCamera() {
  if (userMedia.value) {
    userMedia.value.getVideoTracks().forEach((track) => {
      userMedia.value?.removeTrack(track);
      userMediaProperties.video = false;
    });
  }
}

export async function openScreen() {
  if (!displayMedia.value) {
    displayMedia.value = await navigator.mediaDevices.getDisplayMedia();
    displayMediaActive.value = displayMedia.value.getTracks().length > 0;
  } else {
    const media = await navigator.mediaDevices.getDisplayMedia();
    if (media) {
      media.getTracks().forEach((track) => {
        displayMedia.value?.addTrack(track);
        displayMediaActive.value = true;
      });
    }
  }
}
export function closeScreen() {
  if (displayMedia.value) {
    displayMedia.value.getTracks().forEach((track) => {
      displayMedia.value?.removeTrack(track);
      displayMediaActive.value = false;
    });
    displayMedia.value = undefined;
  }
}

export function closeAllMedia() {
  if (userMedia.value) {
    const tracks = userMedia.value.getTracks();
    tracks.forEach((track) => {
      userMedia.value?.removeTrack(track);
    });
    userMedia.value = undefined;
    userMediaProperties.audio = false;
    userMediaProperties.video = false;
  }
  if (displayMedia.value) {
    const tracks = displayMedia.value.getTracks();
    tracks.forEach((track) => {
      track.stop();
    });
    displayMedia.value = undefined;
    displayMediaActive.value = false;
  }
}
