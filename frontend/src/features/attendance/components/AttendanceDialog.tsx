import { useState, useRef, useCallback, useMemo } from "react";
import Webcam from "react-webcam";
import {
  MapPin,
  Camera,
  RefreshCw,
  CheckCircle2,
  AlertCircle,
} from "lucide-react";
import { toast } from "sonner";

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
    [onDragEnd]
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
    null
  );
  const [isManualLocation, setIsManualLocation] = useState(false);
  const [isLoadingLocation, setIsLoadingLocation] = useState(false);

  const { mutate: clock, isPending } = useClock();

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
        "Failed to load map. Please ensure you have a stable internet connection."
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
        { enableHighAccuracy: true, timeout: 5000, maximumAge: 0 }
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

    clock(
      {
        latitude: location.lat,
        longitude: location.lng,
        image_base64: imgSrc,
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
      }
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
            {step === "scan" ? (
              <Webcam
                audio={false}
                ref={webcamRef}
                screenshotFormat="image/jpeg"
                videoConstraints={{ facingMode: "user" }}
                className="w-full h-full object-cover"
                onUserMediaError={() => setErrorMsg("Camera permission denied")}
              />
            ) : (
              <img
                src={imgSrc!}
                alt="Attendance Preview"
                className="w-full h-full object-cover"
              />
            )}

            {isPending && (
              <div className="absolute inset-0 bg-black/50 flex items-center justify-center text-white backdrop-blur-sm">
                <RefreshCw className="animate-spin h-8 w-8" />
              </div>
            )}
          </div>

          {errorMsg && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{errorMsg}</AlertDescription>
            </Alert>
          )}

          {isLoadingLocation && (
            <div className="text-center text-sm animate-pulse">
              Detecting Location...
            </div>
          )}

          {!isLoadingLocation && location && isManualLocation && (
            <div className="h-48 w-full rounded-md overflow-hidden border border-amber-300 relative">
              <div className="absolute top-0 left-0 z-[1000] bg-amber-100 text-amber-800 text-xs px-2 py-1 rounded-b mx-auto left-0 right-0 w-fit font-medium shadow-sm">
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

          {!isLoadingLocation && location && !isManualLocation && (
            <div className="flex items-center justify-center text-green-600 bg-green-50 p-2 rounded text-sm">
              <MapPin className="w-4 h-4 mr-2" />
              GPS Locked: {location.lat.toFixed(5)}, {location.lng.toFixed(5)}
            </div>
          )}

          <div className="flex gap-3 w-full mt-2">
            {step === "scan" ? (
              <Button onClick={capture} className="w-full" size="lg">
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
