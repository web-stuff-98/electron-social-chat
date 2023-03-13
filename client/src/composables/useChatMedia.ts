import { onMounted, reactive, ref, Ref, watch } from "vue";

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
    try {
      userMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: options.value.userMedia.video,
      });
      trackIds.userMediaVideo = userMediaStream.getVideoTracks()[0].id;
      trackIds.userMediaAudio = userMediaStream.getAudioTracks()[0].id;
    } catch (e) {
      console.warn(e);
    }
    try {
      displayMediaStream = await navigator.mediaDevices.getDisplayMedia({
        audio: true,
        video: options.value.displayMedia.video,
      });
      trackIds.displayMediaVideo = displayMediaStream.getVideoTracks()[0].id;
      trackIds.displayMediaAudio = displayMediaStream.getAudioTracks()[0].id;
    } catch (e) {
      console.warn(e);
    }
    userMediaStream?.getAudioTracks().forEach((track) => {
      track.enabled = options.value.userMedia.audio;
    });
    displayMediaStream?.getAudioTracks().forEach((track) => {
      track.enabled = options.value.displayMedia.audio;
    });
    negotiateConnection();
  });

  watch(options, async (oldVal, newVal) => {
    if (
      newVal.userMedia.video !== oldVal.userMedia.video ||
      newVal.displayMedia.video !== newVal.displayMedia.video
    ) {
      let userMediaStream: MediaStream | undefined;
      let displayMediaStream: MediaStream | undefined;
      try {
        userMediaStream = await navigator.mediaDevices.getUserMedia({
          audio: true,
          video: newVal.userMedia.video,
        });
      } catch (e) {
        console.warn(e);
      }
      try {
        displayMediaStream = await navigator.mediaDevices.getDisplayMedia({
          audio: true,
          video: newVal.displayMedia.video,
        });
      } catch (e) {
        console.warn(e);
      }
      userMediaStream?.getAudioTracks().forEach((track) => {
        track.enabled = options.value.userMedia.audio;
      });
      displayMediaStream?.getAudioTracks().forEach((track) => {
        track.enabled = options.value.displayMedia.audio;
      });
      negotiateConnection();
    } else {
      if (
        newVal.userMedia.audio !== oldVal.userMedia.audio ||
        newVal.displayMedia.audio !== newVal.displayMedia.audio
      ) {
        if (stream.value) {
          stream.value.getTrackById(trackIds.displayMediaAudio)!.enabled =
            options.value.displayMedia.audio;
          stream.value.getTrackById(trackIds.userMediaAudio)!.enabled =
            options.value.userMedia.audio;
        }
      }
    }
  });

  return {
    stream,
    trackIds,
  };
};
