import { onBeforeUnmount, onMounted, reactive, ref, Ref, watch } from "vue";

export const useChatMedia = (
  negotiateConnection: Function,
  options: Ref<{
    userMedia: {
      audio: boolean;
      video: boolean;
    };
    displayMedia: {
      audio: boolean;
      video: boolean;
    };
  }>
) => {
  const stream = ref<MediaStream | undefined>();
  const trackIds = reactive({
    userMediaVideo: "",
    userMediaAudio: "",
    displayMediaVideo: "",
    displayMediaAudio: "",
  });

  onMounted(async () => {
    let userMediaStream: MediaStream | undefined;
    let displayMediaStream: MediaStream | undefined;
    stream.value = new MediaStream();
    try {
      userMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: options.value.userMedia.audio,
        // has to be true or it throws an error.
        video: true,
      });
      const vidTrack = userMediaStream.getVideoTracks()[0];
      const sndTrack = userMediaStream.getAudioTracks()[0];
      trackIds.userMediaVideo = vidTrack !== undefined ? vidTrack.id : "";
      trackIds.userMediaAudio = sndTrack !== undefined ? sndTrack.id : "";
    } catch (e) {
      console.warn(e);
    }
    if (options.value.displayMedia.video) {
      if (options.value.displayMedia.video) {
        try {
          displayMediaStream = await navigator.mediaDevices.getDisplayMedia({
            audio: options.value.displayMedia.audio,
            // has to be true or it throws an error.
            video: true,
          });
          trackIds.displayMediaVideo =
            displayMediaStream.getVideoTracks()[0].id || "";
          trackIds.displayMediaAudio =
            displayMediaStream.getAudioTracks()[0].id || "";
        } catch (e) {
          console.warn(e);
        }
      }
    } else {
      trackIds.displayMediaVideo = "";
      trackIds.displayMediaAudio = "";
    }
    userMediaStream?.getTracks().forEach((track) => {
      if (track.kind === "video" && !options.value.userMedia.video) {
        track.enabled = false;
      }
      stream.value?.addTrack(track);
    });
    displayMediaStream?.getTracks().forEach((track) => {
      if (track.kind === "video" && !options.value.displayMedia.video) {
        track.enabled = false;
      }
      stream.value?.addTrack(track);
    });
    negotiateConnection();
  });

  onBeforeUnmount(() => {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => {
        stream.value?.removeTrack(track);
      });
    }
  });

  watch(options.value, async () => {
    let userMediaStream: MediaStream | undefined;
    let displayMediaStream: MediaStream | undefined;
    stream.value?.getTracks().forEach((track) => {
      stream.value?.removeTrack(track);
    });
    stream.value = new MediaStream();
    try {
      userMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: options.value.userMedia.audio,
        // has to be true or it throws an error.
        video: true,
      });
      const vidTrack = userMediaStream.getVideoTracks()[0];
      const sndTrack = userMediaStream.getAudioTracks()[0];
      trackIds.userMediaVideo = vidTrack !== undefined ? vidTrack.id : "";
      trackIds.userMediaAudio = sndTrack !== undefined ? sndTrack.id : "";
    } catch (e) {
      console.warn(e);
    }
    if (options.value.displayMedia.video) {
      try {
        displayMediaStream = await navigator.mediaDevices.getDisplayMedia({
          audio: options.value.displayMedia.audio,
          // has to be true or it throws an error.
          video: true,
        });
        trackIds.displayMediaVideo =
          displayMediaStream.getVideoTracks()[0].id || "";
        trackIds.displayMediaAudio =
          displayMediaStream.getAudioTracks()[0].id || "";
      } catch (e) {
        console.warn(e);
      }
    } else {
      trackIds.displayMediaVideo = "";
      trackIds.displayMediaAudio = "";
    }
    userMediaStream?.getTracks().forEach((track) => {
      if (track.kind === "video" && !options.value.userMedia.video) {
        track.enabled = false;
      }
      stream.value?.addTrack(track);
    });
    displayMediaStream?.getTracks().forEach((track) => {
      if (track.kind === "video" && !options.value.displayMedia.video) {
        track.enabled = false;
      }
      stream.value?.addTrack(track);
    });
    negotiateConnection();
  });

  return {
    stream,
    trackIds,
  };
};
