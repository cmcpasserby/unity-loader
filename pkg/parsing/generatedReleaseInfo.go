package parsing

import (
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"path/filepath"
	"strings"
)

const unityPathFragment = "{UNITY_PATH}"

var firstVersionWithDocsCategory = unity.VersionFromString("2018.2.0a1")

func getDestination(componentId string) string {
	switch componentId {
	case "mono":
	case "visualstudio":
		return ""
	case "monodevelop":
	case "documentation":
		return unityPathFragment
	case "standardassets":
		return filepath.Join(unityPathFragment, "Standard Assets")
	case "exampleprojects":
	case "example":
		return "/Users/Shared/Unity"
	case "android":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer")
	case "android-sdk-build-tools":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK/build-tools")
	case "android-sdk-platforms":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK/platforms")
	case "android-sdk-platform-tools":
	case "android-sdk-ndk-tools":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK")
	case "android-ndk":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/NDK")
	case "android-open-jdk":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/OpenJDK")
	case "ios":
		return filepath.Join(unityPathFragment, "PlaybackEngines")
	case "tvos":
	case "appletv":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AppleTVSupport")
	case "linux":
		return filepath.Join(unityPathFragment, "PlaybackEngines/LinuxStandaloneSupport")
	case "mac":
	case "mac-il2cpp":
		return filepath.Join(unityPathFragment, "Unity.app/Contents/PlaybackEngines/MacStandaloneSupport")
	case "samsungtv":
	case "samsung-tv":
		return filepath.Join(unityPathFragment, "PlaybackEngines/STVPlayer")
	case "tizen":
		return filepath.Join(unityPathFragment, "PlaybackEngines/TizenPlayer")
	case "vuforia":
	case "vuforia-ar":
		return filepath.Join(unityPathFragment, "PlaybackEngines/VuforiaSupport")
	case "webgl":
		return filepath.Join(unityPathFragment, "PlaybackEngines/WebGLSupport")
	case "windows":
	case "windows-mono":
		return filepath.Join(unityPathFragment, "PlaybackEngines/WindowsStandaloneSupport")
	case "facebook":
	case "facebook-games":
		return filepath.Join(unityPathFragment, "PlaybackEngines/Facebook")
	case "facebookgameroom":
		return ""
	case "lumin":
		return filepath.Join(unityPathFragment, "PlaybackEngines/LuminSupport")
	}

	if strings.HasPrefix(componentId, "language-") {
		return filepath.Join(unityPathFragment, "Unity.app/Contents/Localization")
	}

	return unityPathFragment
}

func getCategory(componentId string, version *CacheVersion) string {
	switch componentId {
	case "monodevelop":
	case "visualstudio":
		return "Dev tools"
	case "mono":
	case "visualstudioprofessionalunityworkload":
	case "visualstudioenterpriseunityworkload":
	case "facebookgameroom":
		return "Plugins"

	case "standardassets":
	case "exampleprojects":
	case "example":
		return "Components"

	case "documentation":
		if version.Compare(firstVersionWithDocsCategory) >= 0 {
			return "Documentation"
		}
		return "Components"
	}

	if strings.HasPrefix(componentId, "language-") {
		return "Language Packs (Preview)"
	}
	return "Platforms"
}

type androidSdkDownloadInfo struct {
	Url           string
	Version       string
	Main          bool
	InstalledSize float32 // MB
	DownloadSize  float32 // MB
	RenameFrom    string
	RenameTo      string
}

// 2019.1:
// Android SDK Tools 26.1.1
// Android SDK Platform tools 28.0.1
// Android SDK Build tools 8.0.2
// Android SDK Platform2 28
// Android NDK r16b
var androidNDKSDKDownloadInfo = map[string]androidSdkDownloadInfo{
	"Android SDK & NDK Tools": {
		Url:           "https://dl.google.com/android/repository/sdk-tools-darwin-4333796.zip",
		Version:       "26.1.1",
		Main:          true,
		InstalledSize: 174, // MB
		DownloadSize:  148, // MB
	},
	"Android SDK Platform Tools": {
		Url:           "https://dl.google.com/android/repository/platform-tools_r28.0.1-darwin.zip",
		Version:       "28.0.1",
		InstalledSize: 15.7, // MB
		DownloadSize:  4.55, // MB
	},
	"Android SDK Build Tools": {
		Url:           "https://dl.google.com/android/repository/build-tools_r28.0.3-macosx.zip",
		Version:       "28.0.3",
		InstalledSize: 120,  // MB
		DownloadSize:  52.6, // MB
		RenameFrom:    "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/SDK/build-tools/android-9",
		RenameTo:      "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/SDK/build-tools/28.0.3",
	},
	"Android SDK Platforms": {
		Url:           "https://dl.google.com/android/repository/platform-28_r06.zip",
		Version:       "28",
		InstalledSize: 121,  // MB
		DownloadSize:  60.6, // MB
		RenameFrom:    "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/SDK/platforms/android-9",
		RenameTo:      "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/SDK/platforms/android-28",
	},
	"Android NDK": {
		Url:           "https://dl.google.com/android/repository/android-ndk-r19b-darwin-x86_64.zip",
		Version:       "r19b",
		InstalledSize: 2700, // MB
		DownloadSize:  770,  // MB
		RenameFrom:    "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/NDK/android-ndk-r19b",
		RenameTo:      "{UNITY_PATH}/PlaybackEngines/AndroidPlayer/NDK",
	},
}

// 2019.2
// Android JDK 8u172-b11
var androidOpenJdkDownloadInfo = androidSdkDownloadInfo{
	Url:           "http://download.unity3d.com/download_unity/open-jdk/open-jdk-mac-x64/jdk8u172-b11_4be8440cc514099cfe1b50cbc74128f6955cd90fd5afe15ea7be60f832de67b4.zip",
	Version:       "8u172-b11",
	Main:          true,
	InstalledSize: 72.7, // MB
	DownloadSize:  165,  // MB
}
