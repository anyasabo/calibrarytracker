declare module 'leaflet.markercluster' {
	import * as L from 'leaflet';

	export class MarkerClusterGroup extends L.FeatureGroup {
		constructor(options?: MarkerClusterGroupOptions);
		addLayer(layer: L.Layer): this;
		removeLayer(layer: L.Layer): this;
		clearLayers(): this;
		addLayers(layers: L.Layer[]): this;
		removeLayers(layers: L.Layer[]): this;
	}

	export interface MarkerClusterGroupOptions extends L.LayerOptions {
		maxClusterRadius?: number;
		showCoverageOnHover?: boolean;
		zoomToBoundsOnClick?: boolean;
		spiderfyOnMaxZoom?: boolean;
		disableClusteringAtZoom?: number;
		chunkedLoading?: boolean;
		iconCreateFunction?: (cluster: MarkerCluster) => L.Icon | L.DivIcon;
	}

	export interface MarkerCluster extends L.Marker {
		getChildCount(): number;
		getAllChildMarkers(): L.Marker[];
	}
}
