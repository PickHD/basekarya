import { useState, useRef, useCallback, useMemo, useEffect } from "react";
import Webcam from "react-webcam";
import {
  MapPin,
  Camera,
  RefreshCw,
  CheckCircle2,
  AlertCircle,
  Loader2,
} from "lucide-react";
import { toast } from "sonner";
import * as faceapi from "face-api.js";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

import { MapContainer, TileLayer, Marker } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useClock } from "@/features/attendance/hooks/useAttendance";

const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});
L.Marker.prototype.options.icon = DefaultIcon;

function DraggableMarker({
  position,
  onDragEnd,
}: {
  position: { lat: number; lng: number };
  onDragEnd: (pos: { lat: number; lng: number }) => void;
}) {
  const markerRef = useRef<L.Marker | null>(null);
  const eventHandlers = useMemo(
    () => ({
      dragend() {
        const marker = markerRef.current;
        if (marker != null) {
          const latlng = marker.getLatLng();
          onDragEnd({ lat: latlng.lat, lng: latlng.lng });
        }
      },
    }),
    [onDragEnd],
  );

  return (
    <Marker
      draggable={true}
      eventHandlers={eventHandlers}
      position={position}
      ref={markerRef}
    />
  );
}

interface AttendanceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  type: "check-in" | "check-out";
}

export function AttendanceDialog({
  open,
  onOpenChange,
  type,
}: AttendanceDialogProps) {
  const webcamRef = useRef<Webcam>(null);
  const [step, setStep] = useState<"scan" | "preview">("scan");
  const [imgSrc, setImgSrc] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [location, setLocation] = useState<{ lat: number; lng: number } | null>(
    null,
  );
  const [isManualLocation, setIsManualLocation] = useState(false);
  const [isLoadingLocation, setIsLoadingLocation] = useState(false);

  const [isModelLoaded, setIsModelLoaded] = useState(false);
  const [isValidFace, setIsValidFace] = useState(false);
  const [faceError, setFaceError] = useState<string>("");

  const { mutate: clock, isPending } = useClock();

  useEffect(() => {
    if (!open) return;

    const loadModels = async () => {
      try {
        setIsModelLoaded(false);
        await faceapi.nets.tinyFaceDetector.loadFromUri("/models");
        setIsModelLoaded(true);
      } catch (err) {
        console.error("Failed to load face models:", err);
        setFaceError("Gagal memuat AI pendeteksi wajah.");
      }
    };

    loadModels();
  }, [open]);

  useEffect(() => {
    if (!isModelLoaded || step !== "scan" || !open) return;

    const interval = setInterval(async () => {
      if (
        webcamRef.current &&
        webcamRef.current.video &&
        webcamRef.current.video.readyState === 4
      ) {
        const video = webcamRef.current.video;

        const detections = await faceapi.detectAllFaces(
          video,
          new faceapi.TinyFaceDetectorOptions({
            inputSize: 224,
            scoreThreshold: 0.5,
          }),
        );

        if (detections.length === 0) {
          setIsValidFace(false);
          setFaceError("Wajah tidak terdeteksi. Posisikan wajah di kamera.");
        } else if (detections.length > 1) {
          setIsValidFace(false);
          setFaceError("Terdeteksi lebih dari 1 wajah! Mohon absen sendiri.");
        } else {
          setIsValidFace(true);
          setFaceError("");
        }
      }
    }, 500);

    return () => clearInterval(interval);
  }, [isModelLoaded, step, open]);

  useEffect(() => {
    if (!open) {
      setStep("scan");
      setImgSrc(null);
      setErrorMsg(null);
      setFaceError("");
      setIsValidFace(false);
      setLocation(null);
    }
  }, [open]);

  const fallbackToIpAndMap = useCallback(async () => {
    try {
      toast.info("GPS failed. Please mark your location on the map.");
      const res = await fetch("https://ipapi.co/json/");
      const data = await res.json();

      if (data.latitude && data.longitude) {
        setLocation({ lat: data.latitude, lng: data.longitude });
        setIsManualLocation(true);
      }
    } catch (e) {
      console.error(e);
      toast.error(
        "Failed to load map. Please ensure you have a stable internet connection.",
      );
    } finally {
      setIsLoadingLocation(false);
    }
  }, []);

  const getLocation = useCallback(() => {
    setIsLoadingLocation(true);
    setIsManualLocation(false);

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          setLocation({ lat: pos.coords.latitude, lng: pos.coords.longitude });
          setIsLoadingLocation(false);
          toast.success("GPS Accurate Locked!");
        },
        (err) => {
          console.warn("GPS failed, fallback to IP + Manual Map", err);
          fallbackToIpAndMap();
        },
        { enableHighAccuracy: true, timeout: 5000, maximumAge: 0 },
      );
    } else {
      fallbackToIpAndMap();
    }
  }, [fallbackToIpAndMap]);

  const capture = useCallback(() => {
    const imageSrc = webcamRef.current?.getScreenshot();
    if (imageSrc) {
      setImgSrc(imageSrc);
      setStep("preview");
      getLocation();
    }
  }, [webcamRef, getLocation]);

  const retake = () => {
    setImgSrc(null);
    setStep("scan");
    setErrorMsg(null);
  };

  const handleSubmit = () => {
    if (!imgSrc || !location) return;

    const rawBase64 = imgSrc.split(",")[1];

    clock(
      {
        latitude: location.lat,
        longitude: location.lng,
        image_base64: rawBase64,
        notes: isManualLocation
          ? "[MANUAL] User adjusted location on map"
          : "[GPS] Auto-detected",
      },
      {
        onSuccess: () => {
          onOpenChange(false);
          setImgSrc(null);
          setStep("scan");
        },
      },
    );
  };

  const title =
    type === "check-in" ? "Clock In Attendance" : "Clock Out Attendance";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            Please ensure your face is visible and location is enabled.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col items-center gap-4 py-4">
          <div className="relative w-full aspect-[4/3] bg-slate-950 rounded-lg overflow-hidden border-2 border-slate-200 shadow-inner">
            {step === "scan" && !isModelLoaded && (
              <div className="absolute inset-0 z-20 flex flex-col items-center justify-center bg-slate-900/90 text-white">
                <Loader2 className="w-8 h-8 animate-spin mb-2" />
                <span className="text-sm">Memuat AI Wajah...</span>
              </div>
            )}

            {step === "scan" ? (
              <>
                <Webcam
                  audio={false}
                  ref={webcamRef}
                  screenshotFormat="image/jpeg"
                  screenshotQuality={0.8}
                  videoConstraints={{ facingMode: "user" }}
                  className={`w-full h-full object-cover transform scale-x-[-1] transition-opacity duration-300 ${!isValidFace && isModelLoaded ? "opacity-50" : ""}`}
                  onUserMediaError={() =>
                    setErrorMsg("Camera permission denied")
                  }
                />

                {isValidFace && (
                  <div className="absolute inset-4 border-4 border-green-400/50 rounded-lg z-10 animate-pulse pointer-events-none" />
                )}
              </>
            ) : (
              <img
                src={imgSrc!}
                alt="Attendance Preview"
                className="w-full h-full object-cover transform scale-x-[-1]"
              />
            )}

            {isPending && (
              <div className="absolute inset-0 z-30 bg-black/50 flex items-center justify-center text-white backdrop-blur-sm">
                <RefreshCw className="animate-spin h-8 w-8" />
              </div>
            )}
          </div>

          {errorMsg && (
            <Alert variant="destructive" className="w-full">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{errorMsg}</AlertDescription>
            </Alert>
          )}

          {step === "scan" &&
            !errorMsg &&
            isModelLoaded &&
            (faceError ? (
              <Alert
                variant="destructive"
                className="w-full py-2 flex items-center gap-2 [&>svg]:static [&>svg~*]:pl-0"
              >
                <AlertCircle className="h-4 w-4 shrink-0" />
                <AlertDescription className="text-xs font-medium">
                  {faceError}
                </AlertDescription>
              </Alert>
            ) : (
              <Alert className="w-full py-2 border-green-200 bg-green-50 text-green-700 flex items-center gap-2 [&>svg]:static [&>svg~*]:pl-0">
                <CheckCircle2 className="h-4 w-4 text-green-600 shrink-0" />
                <AlertDescription className="text-xs font-medium">
                  Wajah terverifikasi. Siap absen.
                </AlertDescription>
              </Alert>
            ))}

          {isLoadingLocation && (
            <div className="text-center text-sm animate-pulse">
              Detecting Location...
            </div>
          )}

          {!isLoadingLocation &&
            location &&
            isManualLocation &&
            step === "preview" && (
              <div className="h-48 w-full rounded-md overflow-hidden border border-amber-300 relative">
                <div className="absolute top-0 left-0 z-[1000] bg-amber-100 text-amber-800 text-xs px-2 py-1 rounded-b mx-auto right-0 w-fit font-medium shadow-sm">
                  GPS weak. Drag the pin to your current location.
                </div>
                <MapContainer
                  center={location}
                  zoom={15}
                  scrollWheelZoom={false}
                  style={{ height: "100%", width: "100%", zIndex: 1 }}
                >
                  <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                  <DraggableMarker
                    position={location}
                    onDragEnd={(newPos) => setLocation(newPos)}
                  />
                </MapContainer>
              </div>
            )}

          {!isLoadingLocation &&
            location &&
            !isManualLocation &&
            step === "preview" && (
              <div className="flex items-center justify-center text-green-600 bg-green-50 p-2 rounded text-sm w-full">
                <MapPin className="w-4 h-4 mr-2" />
                GPS Locked: {location.lat.toFixed(5)}, {location.lng.toFixed(5)}
              </div>
            )}

          <div className="flex gap-3 w-full mt-2">
            {step === "scan" ? (
              <Button
                onClick={capture}
                className="w-full"
                size="lg"
                disabled={!isValidFace || !isModelLoaded || !!errorMsg}
              >
                <Camera className="mr-2 h-4 w-4" /> Capture Photo
              </Button>
            ) : (
              <div className="flex w-full gap-2">
                <Button
                  variant="outline"
                  onClick={retake}
                  disabled={isPending}
                  className="flex-1"
                >
                  Retake
                </Button>
                <Button
                  onClick={handleSubmit}
                  disabled={isPending || !location}
                  className="w-full"
                >
                  {isPending ? (
                    <>
                      <RefreshCw className="mr-2 h-4 w-4 animate-spin" />{" "}
                      Processing...
                    </>
                  ) : (
                    <>
                      <CheckCircle2 className="mr-2 h-4 w-4" /> Confirm{" "}
                      {type === "check-in" ? "In" : "Out"}
                    </>
                  )}
                </Button>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
