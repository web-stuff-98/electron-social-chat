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
  const userStream = ref<MediaStream | undefined>();
  const displayStream = ref<MediaStream | undefined>();
  const streamIds = reactive({
    userMedia: "",
    displayMedia: "",
  });

  onMounted(async () => {
    let userMediaStream: MediaStream | undefined;
    let displayMediaStream: MediaStream | undefined;
    userStream.value = new MediaStream();
    displayStream.value = new MediaStream();
    try {
      userMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: options.value.userMedia.audio
          ? {
              noiseSuppression: true,
              echoCancellation: true,
            }
          : false,
        // has to be true or it throws an error.
        video: true,
      });
      const vidTrack = userMediaStream.getVideoTracks()[0];
      const sndTrack = userMediaStream.getAudioTracks()[0];
      if (!options.value.userMedia.video) {
        if (vidTrack !== undefined) {
          vidTrack.enabled = false;
        }
      } else {
        vidTrack.contentHint = "motion";
      }
      if (sndTrack) {
        sndTrack.contentHint = "speech";
        streamIds.userMedia = userMediaStream.id;
      }
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
        const vidTrack = displayMediaStream.getVideoTracks()[0];
        if (vidTrack) {
          streamIds.displayMedia = vidTrack.id;
          vidTrack.contentHint = "detail";
          const sndTrack = displayMediaStream.getAudioTracks()[0];
          if (sndTrack) {
            sndTrack.contentHint = "music";
          }
        } else {
          streamIds.displayMedia = "";
        }
      } catch (e) {
        console.warn(e);
      }
    }
    userStream.value = userMediaStream;
    displayStream.value = displayMediaStream;
    negotiateConnection(true);
  });

  onBeforeUnmount(() => {
    if (userStream.value) {
      userStream.value.getTracks().forEach((track) => {
        userStream.value?.removeTrack(track);
      });
    }
    if (displayStream.value) {
      displayStream.value.getTracks().forEach((track) => {
        displayStream.value?.removeTrack(track);
      });
    }
  });

  watch(options.value, async () => {
    let userMediaStream: MediaStream | undefined;
    let displayMediaStream: MediaStream | undefined;
    if (userStream.value) {
      userStream.value?.getTracks().forEach((track) => {
        userStream.value?.removeTrack(track);
      });
    }
    if (displayStream.value) {
      displayStream.value?.getTracks().forEach((track) => {
        displayStream.value?.removeTrack(track);
      });
    }
    userStream.value = new MediaStream();
    displayStream.value = new MediaStream();
    try {
      userMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: options.value.userMedia.audio
          ? {
              noiseSuppression: true,
              echoCancellation: true,
            }
          : false,
        // has to be true or it throws an error.
        video: true,
      });
      const vidTrack = userMediaStream.getVideoTracks()[0];
      const sndTrack = userMediaStream.getAudioTracks()[0];
      if (!options.value.userMedia.video) {
        if (vidTrack !== undefined) {
          vidTrack.enabled = false;
        }
      } else {
        vidTrack.contentHint = "motion";
      }
      if (sndTrack) {
        sndTrack.contentHint = "speech";
        streamIds.userMedia = userMediaStream.id;
      }
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
        const vidTrack = displayMediaStream.getVideoTracks()[0];
        if (vidTrack) {
          streamIds.displayMedia = vidTrack.id;
          vidTrack.contentHint = "detail";
        } else {
          streamIds.displayMedia = "";
        }
        const sndTrack = displayMediaStream.getAudioTracks()[0];
        if (sndTrack) {
          sndTrack.contentHint = "music";
        }
      } catch (e) {
        console.warn(e);
      }
    }
    userStream.value = userMediaStream;
    displayStream.value = displayMediaStream;
    negotiateConnection();
  });

  return {
    userStream,
    displayStream,
    streamIds,
  };
};
