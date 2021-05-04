<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.4" tiledversion="1.4.3" name="objects" tilewidth="16" tileheight="16" tilecount="270" columns="16" objectalignment="topleft">
 <grid orientation="orthogonal" width="1" height="1"/>
 <tile id="0">
  <image width="16" height="16" source="nosprite.png"/>
 </tile>
 <tile id="1" type="WarpZone">
  <properties>
   <property name="orientation" value="EN"/>
  </properties>
  <image width="16" height="16" source="warpzone_en.png"/>
 </tile>
 <tile id="2" type="WarpZone">
  <properties>
   <property name="orientation" value="ES"/>
  </properties>
  <image width="16" height="16" source="warpzone_es.png"/>
 </tile>
 <tile id="3" type="WarpZone">
  <properties>
   <property name="orientation" value="NE"/>
  </properties>
  <image width="16" height="16" source="warpzone_ne.png"/>
 </tile>
 <tile id="4" type="WarpZone">
  <properties>
   <property name="orientation" value="NW"/>
  </properties>
  <image width="16" height="16" source="warpzone_nw.png"/>
 </tile>
 <tile id="5" type="WarpZone">
  <properties>
   <property name="orientation" value="SE"/>
  </properties>
  <image width="16" height="16" source="warpzone_se.png"/>
 </tile>
 <tile id="6" type="WarpZone">
  <properties>
   <property name="orientation" value="SW"/>
  </properties>
  <image width="16" height="16" source="warpzone_sw.png"/>
 </tile>
 <tile id="7" type="WarpZone">
  <properties>
   <property name="orientation" value="WN"/>
  </properties>
  <image width="16" height="16" source="warpzone_wn.png"/>
 </tile>
 <tile id="8" type="WarpZone">
  <properties>
   <property name="orientation" value="WS"/>
  </properties>
  <image width="16" height="16" source="warpzone_ws.png"/>
 </tile>
 <tile id="9" type="Sprite">
  <properties>
   <property name="image_dir" value="sprites"/>
   <property name="player_solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="playerclip.png"/>
 </tile>
 <tile id="10" type="Sprite">
  <properties>
   <property name="image_dir" value="sprites"/>
   <property name="object_solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="objectclip.png"/>
 </tile>
 <tile id="11" type="RiserFsck">
  <image width="16" height="16" source="riserfsck.png"/>
 </tile>
 <tile id="12" type="OneWay">
  <properties>
   <property name="orientation" value="EN"/>
  </properties>
  <image width="16" height="16" source="editorimgs/oneway_e.png"/>
 </tile>
 <tile id="13" type="OneWay">
  <properties>
   <property name="orientation" value="NW"/>
  </properties>
  <image width="16" height="16" source="editorimgs/oneway_n.png"/>
 </tile>
 <tile id="14" type="OneWay">
  <properties>
   <property name="orientation" value="WS"/>
  </properties>
  <image width="16" height="16" source="editorimgs/oneway_w.png"/>
 </tile>
 <tile id="15" type="OneWay">
  <properties>
   <property name="orientation" value="SE"/>
  </properties>
  <image width="16" height="16" source="editorimgs/oneway_s.png"/>
 </tile>
</tileset>
