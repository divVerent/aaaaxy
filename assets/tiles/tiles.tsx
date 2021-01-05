<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.4" tiledversion="1.4.3" name="tiles" tilewidth="16" tileheight="16" tilecount="56" columns="0" objectalignment="topleft">
 <grid orientation="orthogonal" width="1" height="1"/>
 <tile id="0">
  <properties>
   <property name="opaque" value="false"/>
   <property name="solid" value="false"/>
  </properties>
  <image width="16" height="16" source="empty.png"/>
 </tile>
 <tile id="1" type="WarpZone">
  <properties>
   <property name="orientation" value="EN"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_en.png"/>
 </tile>
 <tile id="2" type="WarpZone">
  <properties>
   <property name="orientation" value="ES"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_es.png"/>
 </tile>
 <tile id="3" type="WarpZone">
  <properties>
   <property name="orientation" value="NE"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_ne.png"/>
 </tile>
 <tile id="4" type="WarpZone">
  <properties>
   <property name="orientation" value="NW"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_nw.png"/>
 </tile>
 <tile id="5" type="WarpZone">
  <properties>
   <property name="orientation" value="SE"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_se.png"/>
 </tile>
 <tile id="6" type="WarpZone">
  <properties>
   <property name="orientation" value="SW"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_sw.png"/>
 </tile>
 <tile id="7" type="WarpZone">
  <properties>
   <property name="orientation" value="WN"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_wn.png"/>
 </tile>
 <tile id="8" type="WarpZone">
  <properties>
   <property name="orientation" value="WS"/>
   <property name="type" value="WarpZone"/>
  </properties>
  <image width="16" height="16" source="warpzone_ws.png"/>
 </tile>
 <!--
 <tile id="9">
  <properties>
   <property name="img.EN" value="wall_.png"/>
   <property name="img.ES" value="wall_.png"/>
   <property name="img.NE" value="wall_.png"/>
   <property name="img.NW" value="wall_.png"/>
   <property name="img.SE" value="wall_.png"/>
   <property name="img.SW" value="wall_.png"/>
   <property name="img.WN" value="wall_.png"/>
   <property name="img.WS" value="wall_.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_.png"/>
 </tile>
 <tile id="10">
  <properties>
   <property name="img.EN" value="wall_w.png"/>
   <property name="img.ES" value="wall_w.png"/>
   <property name="img.NE" value="wall_s.png"/>
   <property name="img.NW" value="wall_s.png"/>
   <property name="img.SE" value="wall_n.png"/>
   <property name="img.SW" value="wall_n.png"/>
   <property name="img.WN" value="wall_e.png"/>
   <property name="img.WS" value="wall_e.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_w.png"/>
 </tile>
 <tile id="11">
  <properties>
   <property name="img.EN" value="wall_s.png"/>
   <property name="img.ES" value="wall_n.png"/>
   <property name="img.NE" value="wall_w.png"/>
   <property name="img.NW" value="wall_e.png"/>
   <property name="img.SE" value="wall_w.png"/>
   <property name="img.SW" value="wall_e.png"/>
   <property name="img.WN" value="wall_s.png"/>
   <property name="img.WS" value="wall_n.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s.png"/>
 </tile>
 <tile id="12">
  <properties>
   <property name="img.EN" value="wall_s_w_sw0.png"/>
   <property name="img.ES" value="wall_n_w_nw0.png"/>
   <property name="img.NE" value="wall_s_w_sw0.png"/>
   <property name="img.NW" value="wall_e_s_se0.png"/>
   <property name="img.SE" value="wall_n_w_nw0.png"/>
   <property name="img.SW" value="wall_n_e_ne0.png"/>
   <property name="img.WN" value="wall_e_s_se0.png"/>
   <property name="img.WS" value="wall_n_e_ne0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s_w_sw0.png"/>
 </tile>
 <tile id="13">
  <properties>
   <property name="img.EN" value="wall_s_w_sw1.png"/>
   <property name="img.ES" value="wall_n_w_nw1.png"/>
   <property name="img.NE" value="wall_s_w_sw1.png"/>
   <property name="img.NW" value="wall_e_s_se1.png"/>
   <property name="img.SE" value="wall_n_w_nw1.png"/>
   <property name="img.SW" value="wall_n_e_ne1.png"/>
   <property name="img.WN" value="wall_e_s_se1.png"/>
   <property name="img.WS" value="wall_n_e_ne1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_s_w_sw1.png"/>
 </tile>
 <tile id="14">
  <properties>
   <property name="img.EN" value="wall_e.png"/>
   <property name="img.ES" value="wall_e.png"/>
   <property name="img.NE" value="wall_n.png"/>
   <property name="img.NW" value="wall_n.png"/>
   <property name="img.SE" value="wall_s.png"/>
   <property name="img.SW" value="wall_s.png"/>
   <property name="img.WN" value="wall_w.png"/>
   <property name="img.WS" value="wall_w.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e.png"/>
 </tile>
 <tile id="15">
  <properties>
   <property name="img.EN" value="wall_e_w.png"/>
   <property name="img.ES" value="wall_e_w.png"/>
   <property name="img.NE" value="wall_n_s.png"/>
   <property name="img.NW" value="wall_n_s.png"/>
   <property name="img.SE" value="wall_n_s.png"/>
   <property name="img.SW" value="wall_n_s.png"/>
   <property name="img.WN" value="wall_e_w.png"/>
   <property name="img.WS" value="wall_e_w.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_w.png"/>
 </tile>
 <tile id="16">
  <properties>
   <property name="img.EN" value="wall_e_s_se0.png"/>
   <property name="img.ES" value="wall_n_e_ne0.png"/>
   <property name="img.NE" value="wall_n_w_nw0.png"/>
   <property name="img.NW" value="wall_n_e_ne0.png"/>
   <property name="img.SE" value="wall_s_w_sw0.png"/>
   <property name="img.SW" value="wall_e_s_se0.png"/>
   <property name="img.WN" value="wall_s_w_sw0.png"/>
   <property name="img.WS" value="wall_n_w_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_se0.png"/>
 </tile>
 <tile id="17">
  <properties>
   <property name="img.EN" value="wall_e_s_se1.png"/>
   <property name="img.ES" value="wall_n_e_ne1.png"/>
   <property name="img.NE" value="wall_n_w_nw1.png"/>
   <property name="img.NW" value="wall_n_e_ne1.png"/>
   <property name="img.SE" value="wall_s_w_sw1.png"/>
   <property name="img.SW" value="wall_e_s_se1.png"/>
   <property name="img.WN" value="wall_s_w_sw1.png"/>
   <property name="img.WS" value="wall_n_w_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_se1.png"/>
 </tile>
 <tile id="18">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.ES" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.NE" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.SE" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.WN" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.WS" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se0_sw0.png"/>
 </tile>
 <tile id="19">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.ES" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.NE" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.SE" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.WN" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.WS" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se0_sw1.png"/>
 </tile>
 <tile id="20">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.ES" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.NE" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.SE" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.WN" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.WS" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se1_sw0.png"/>
 </tile>
 <tile id="21">
  <properties>
   <property name="img.EN" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.ES" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.NE" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.SE" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.WN" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.WS" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_e_s_w_se1_sw1.png"/>
 </tile>
 <tile id="22">
  <properties>
   <property name="img.EN" value="wall_n.png"/>
   <property name="img.ES" value="wall_s.png"/>
   <property name="img.NE" value="wall_e.png"/>
   <property name="img.NW" value="wall_w.png"/>
   <property name="img.SE" value="wall_e.png"/>
   <property name="img.SW" value="wall_w.png"/>
   <property name="img.WN" value="wall_n.png"/>
   <property name="img.WS" value="wall_s.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n.png"/>
 </tile>
 <tile id="23">
  <properties>
   <property name="img.EN" value="wall_n_w_nw0.png"/>
   <property name="img.ES" value="wall_s_w_sw0.png"/>
   <property name="img.NE" value="wall_e_s_se0.png"/>
   <property name="img.NW" value="wall_s_w_sw0.png"/>
   <property name="img.SE" value="wall_n_e_ne0.png"/>
   <property name="img.SW" value="wall_n_w_nw0.png"/>
   <property name="img.WN" value="wall_n_e_ne0.png"/>
   <property name="img.WS" value="wall_e_s_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_w_nw0.png"/>
 </tile>
 <tile id="24">
  <properties>
   <property name="img.EN" value="wall_n_w_nw1.png"/>
   <property name="img.ES" value="wall_s_w_sw1.png"/>
   <property name="img.NE" value="wall_e_s_se1.png"/>
   <property name="img.NW" value="wall_s_w_sw1.png"/>
   <property name="img.SE" value="wall_n_e_ne1.png"/>
   <property name="img.SW" value="wall_n_w_nw1.png"/>
   <property name="img.WN" value="wall_n_e_ne1.png"/>
   <property name="img.WS" value="wall_e_s_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_w_nw1.png"/>
 </tile>
 <tile id="25">
  <properties>
   <property name="img.EN" value="wall_n_s.png"/>
   <property name="img.ES" value="wall_n_s.png"/>
   <property name="img.NE" value="wall_e_w.png"/>
   <property name="img.NW" value="wall_e_w.png"/>
   <property name="img.SE" value="wall_e_w.png"/>
   <property name="img.SW" value="wall_e_w.png"/>
   <property name="img.WN" value="wall_n_s.png"/>
   <property name="img.WS" value="wall_n_s.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s.png"/>
 </tile>
 <tile id="26">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.NE" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.NW" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.SE" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.WS" value="wall_n_e_s_ne0_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw0_nw0.png"/>
 </tile>
 <tile id="27">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.NE" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.NW" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.SE" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.WS" value="wall_n_e_s_ne0_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw0_nw1.png"/>
 </tile>
 <tile id="28">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.NE" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.NW" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.SE" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.WS" value="wall_n_e_s_ne1_se0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw1_nw0.png"/>
 </tile>
 <tile id="29">
  <properties>
   <property name="img.EN" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.NE" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.NW" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.SE" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.WS" value="wall_n_e_s_ne1_se1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_s_w_sw1_nw1.png"/>
 </tile>
 <tile id="30">
  <properties>
   <property name="img.EN" value="wall_n_e_ne0.png"/>
   <property name="img.ES" value="wall_e_s_se0.png"/>
   <property name="img.NE" value="wall_n_e_ne0.png"/>
   <property name="img.NW" value="wall_n_w_nw0.png"/>
   <property name="img.SE" value="wall_e_s_se0.png"/>
   <property name="img.SW" value="wall_s_w_sw0.png"/>
   <property name="img.WN" value="wall_n_w_nw0.png"/>
   <property name="img.WS" value="wall_s_w_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_ne0.png"/>
 </tile>
 <tile id="31">
  <properties>
   <property name="img.EN" value="wall_n_e_ne1.png"/>
   <property name="img.ES" value="wall_e_s_se1.png"/>
   <property name="img.NE" value="wall_n_e_ne1.png"/>
   <property name="img.NW" value="wall_n_w_nw1.png"/>
   <property name="img.SE" value="wall_e_s_se1.png"/>
   <property name="img.SW" value="wall_s_w_sw1.png"/>
   <property name="img.WN" value="wall_n_w_nw1.png"/>
   <property name="img.WS" value="wall_s_w_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_ne1.png"/>
 </tile>
 <tile id="32">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.ES" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.NE" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.NW" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.SW" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.WS" value="wall_e_s_w_se0_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne0_nw0.png"/>
 </tile>
 <tile id="33">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.ES" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.NE" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.NW" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.SW" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.WS" value="wall_e_s_w_se1_sw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne0_nw1.png"/>
 </tile>
 <tile id="34">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.ES" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.NE" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.NW" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.SW" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.WS" value="wall_e_s_w_se0_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne1_nw0.png"/>
 </tile>
 <tile id="35">
  <properties>
   <property name="img.EN" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.ES" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.NE" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.NW" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.SW" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.WS" value="wall_e_s_w_se1_sw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_w_ne1_nw1.png"/>
 </tile>
 <tile id="36">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.ES" value="wall_n_e_s_ne0_se0.png"/>
   <property name="img.NE" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_w_ne0_nw0.png"/>
   <property name="img.SE" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.SW" value="wall_e_s_w_se0_sw0.png"/>
   <property name="img.WN" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_s_w_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne0_se0.png"/>
 </tile>
 <tile id="37">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.ES" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.NE" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.SE" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.SW" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.WN" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne0_se1.png"/>
 </tile>
 <tile id="38">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne1_se0.png"/>
   <property name="img.ES" value="wall_n_e_s_ne0_se1.png"/>
   <property name="img.NE" value="wall_n_e_w_ne1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_w_ne0_nw1.png"/>
   <property name="img.SE" value="wall_e_s_w_se1_sw0.png"/>
   <property name="img.SW" value="wall_e_s_w_se0_sw1.png"/>
   <property name="img.WN" value="wall_n_s_w_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_s_w_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne1_se0.png"/>
 </tile>
 <tile id="39">
  <properties>
   <property name="img.EN" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.ES" value="wall_n_e_s_ne1_se1.png"/>
   <property name="img.NE" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_w_ne1_nw1.png"/>
   <property name="img.SE" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.SW" value="wall_e_s_w_se1_sw1.png"/>
   <property name="img.WN" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_s_w_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_ne1_se1.png"/>
 </tile>
 <tile id="40">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw0_nw0.png"/>
 </tile>
 <tile id="41">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
 </tile>
 <tile id="42">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
 </tile>
 <tile id="43">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
 </tile>
 <tile id="44">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
 </tile>
 <tile id="45">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
 </tile>
 <tile id="46">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
 </tile>
 <tile id="47">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
 </tile>
 <tile id="48">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw0_nw0.png"/>
 </tile>
 <tile id="49">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
 </tile>
 <tile id="50">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw1_nw0.png"/>
 </tile>
 <tile id="51">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
 </tile>
 <tile id="52">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se0_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw0_nw0.png"/>
 </tile>
 <tile id="53">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
 </tile>
 <tile id="54">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw0_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne0_se1_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se0_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw1_nw0.png"/>
 </tile>
 <tile id="55">
  <properties>
   <property name="img.EN" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.ES" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.NE" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.NW" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.SE" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.SW" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.WN" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="img.WS" value="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
   <property name="opaque" type="bool" value="true"/>
   <property name="solid" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="wall_n_e_s_w_ne1_se1_sw1_nw1.png"/>
 </tile>
 -->

								<tile id="9">
								<image width="16" height="16" source="../tiles/wall_.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_.png" />
								<property name="img.SW" type="string" value="../tiles/wall_.png" />
								<property name="img.WN" type="string" value="../tiles/wall_.png" />
								<property name="img.NE" type="string" value="../tiles/wall_.png" />
								<property name="img.SE" type="string" value="../tiles/wall_.png" />
								<property name="img.EN" type="string" value="../tiles/wall_.png" />
								<property name="img.NW" type="string" value="../tiles/wall_.png" />
								<property name="img.WS" type="string" value="../tiles/wall_.png" />
								</properties>
								</tile>
								<tile id="10">
								<image width="16" height="16" source="../tiles/wall_w.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_w.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e.png" />
								<property name="img.NE" type="string" value="../tiles/wall_s.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n.png" />
								<property name="img.EN" type="string" value="../tiles/wall_w.png" />
								<property name="img.NW" type="string" value="../tiles/wall_s.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e.png" />
								</properties>
								</tile>
								<tile id="11">
								<image width="16" height="16" source="../tiles/wall_s.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_s.png" />
								<property name="img.SW" type="string" value="../tiles/wall_w.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n.png" />
								<property name="img.NW" type="string" value="../tiles/wall_w.png" />
								<property name="img.WS" type="string" value="../tiles/wall_s.png" />
								</properties>
								</tile>
								<tile id="12">
								<image width="16" height="16" source="../tiles/wall_s_w_sw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_se0.png" />
								</properties>
								</tile>
								<tile id="13">
								<image width="16" height="16" source="../tiles/wall_s_w_sw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_se1.png" />
								</properties>
								</tile>
								<tile id="14">
								<image width="16" height="16" source="../tiles/wall_e.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e.png" />
								<property name="img.SW" type="string" value="../tiles/wall_s.png" />
								<property name="img.WN" type="string" value="../tiles/wall_w.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n.png" />
								<property name="img.SE" type="string" value="../tiles/wall_s.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n.png" />
								<property name="img.WS" type="string" value="../tiles/wall_w.png" />
								</properties>
								</tile>
								<tile id="15">
								<image width="16" height="16" source="../tiles/wall_e_w.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_w.png" />
								</properties>
								</tile>
								<tile id="16">
								<image width="16" height="16" source="../tiles/wall_e_s_se0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_s_w_sw0.png" />
								</properties>
								</tile>
								<tile id="17">
								<image width="16" height="16" source="../tiles/wall_e_s_se1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_s_w_sw1.png" />
								</properties>
								</tile>
								<tile id="18">
								<image width="16" height="16" source="../tiles/wall_e_s_w_se0_sw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								</properties>
								</tile>
								<tile id="19">
								<image width="16" height="16" source="../tiles/wall_e_s_w_se0_sw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								</properties>
								</tile>
								<tile id="20">
								<image width="16" height="16" source="../tiles/wall_e_s_w_se1_sw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								</properties>
								</tile>
								<tile id="21">
								<image width="16" height="16" source="../tiles/wall_e_s_w_se1_sw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								</properties>
								</tile>
								<tile id="22">
								<image width="16" height="16" source="../tiles/wall_n.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e.png" />
								<property name="img.WN" type="string" value="../tiles/wall_s.png" />
								<property name="img.NE" type="string" value="../tiles/wall_w.png" />
								<property name="img.SE" type="string" value="../tiles/wall_w.png" />
								<property name="img.EN" type="string" value="../tiles/wall_s.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n.png" />
								</properties>
								</tile>
								<tile id="23">
								<image width="16" height="16" source="../tiles/wall_n_w_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_ne0.png" />
								</properties>
								</tile>
								<tile id="24">
								<image width="16" height="16" source="../tiles/wall_n_w_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_ne1.png" />
								</properties>
								</tile>
								<tile id="25">
								<image width="16" height="16" source="../tiles/wall_n_s.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_s.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_w.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_s.png" />
								</properties>
								</tile>
								<tile id="26">
								<image width="16" height="16" source="../tiles/wall_n_s_w_sw0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								</properties>
								</tile>
								<tile id="27">
								<image width="16" height="16" source="../tiles/wall_n_s_w_sw0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								</properties>
								</tile>
								<tile id="28">
								<image width="16" height="16" source="../tiles/wall_n_s_w_sw1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								</properties>
								</tile>
								<tile id="29">
								<image width="16" height="16" source="../tiles/wall_n_s_w_sw1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								</properties>
								</tile>
								<tile id="30">
								<image width="16" height="16" source="../tiles/wall_n_e_ne0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_w_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_s_w_sw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_se0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_ne0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_w_nw0.png" />
								</properties>
								</tile>
								<tile id="31">
								<image width="16" height="16" source="../tiles/wall_n_e_ne1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_w_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_s_w_sw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_se1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_ne1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_w_nw1.png" />
								</properties>
								</tile>
								<tile id="32">
								<image width="16" height="16" source="../tiles/wall_n_e_w_ne0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								</properties>
								</tile>
								<tile id="33">
								<image width="16" height="16" source="../tiles/wall_n_e_w_ne0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								</properties>
								</tile>
								<tile id="34">
								<image width="16" height="16" source="../tiles/wall_n_e_w_ne1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								</properties>
								</tile>
								<tile id="35">
								<image width="16" height="16" source="../tiles/wall_n_e_w_ne1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								</properties>
								</tile>
								<tile id="36">
								<image width="16" height="16" source="../tiles/wall_n_e_s_ne0_se0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_w_se0_sw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_ne0_se0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_w_ne0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_s_w_sw0_nw0.png" />
								</properties>
								</tile>
								<tile id="37">
								<image width="16" height="16" source="../tiles/wall_n_e_s_ne0_se1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								</properties>
								</tile>
								<tile id="38">
								<image width="16" height="16" source="../tiles/wall_n_e_s_ne1_se0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_ne1_se0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_w_se1_sw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_s_w_sw1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_w_ne0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_w_se0_sw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_ne0_se1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_w_ne1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_s_w_sw0_nw1.png" />
								</properties>
								</tile>
								<tile id="39">
								<image width="16" height="16" source="../tiles/wall_n_e_s_ne1_se1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_e_s_w_se1_sw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_ne1_se1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_w_ne1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_s_w_sw1_nw1.png" />
								</properties>
								</tile>
								<tile id="40">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw0.png" />
								</properties>
								</tile>
								<tile id="41">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								</properties>
								</tile>
								<tile id="42">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								</properties>
								</tile>
								<tile id="43">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								</properties>
								</tile>
								<tile id="44">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								</properties>
								</tile>
								<tile id="45">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								</properties>
								</tile>
								<tile id="46">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								</properties>
								</tile>
								<tile id="47">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								</properties>
								</tile>
								<tile id="48">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw0_nw1.png" />
								</properties>
								</tile>
								<tile id="49">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								</properties>
								</tile>
								<tile id="50">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw0_nw1.png" />
								</properties>
								</tile>
								<tile id="51">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								</properties>
								</tile>
								<tile id="52">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se0_sw1_nw1.png" />
								</properties>
								</tile>
								<tile id="53">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								</properties>
								</tile>
								<tile id="54">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw0.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw0_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se0_sw1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne0_se1_sw1_nw1.png" />
								</properties>
								</tile>
								<tile id="55">
								<image width="16" height="16" source="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<properties>
								<property name="opaque" type="bool" value="true" />
								<property name="solid" type="bool" value="true" />
								<property name="img.ES" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.SW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.WN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.NE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.SE" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.EN" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.NW" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								<property name="img.WS" type="string" value="../tiles/wall_n_e_s_w_ne1_se1_sw1_nw1.png" />
								</properties>
								</tile>
 <wangsets>
  <wangset name="Block/Open" tile="-1">
   <wangedgecolor name="Block" color="#00007f" tile="-1" probability="1"/>
   <wangedgecolor name="Open" color="#aaff7f" tile="-1" probability="1"/>
   <wangcornercolor name="Block" color="#00007f" tile="-1" probability="1"/>
   <wangcornercolor name="Open" color="#aaff7f" tile="-1" probability="1"/>
   <wangtile tileid="9" wangid="0x22222222"/>
   <wangtile tileid="10" wangid="0x21222222"/>
   <wangtile tileid="11" wangid="0x22212222"/>
   <wangtile tileid="12" wangid="0x21212222"/>
   <wangtile tileid="13" wangid="0x21112222"/>
   <wangtile tileid="14" wangid="0x22222122"/>
   <wangtile tileid="15" wangid="0x21222122"/>
   <wangtile tileid="16" wangid="0x22212122"/>
   <wangtile tileid="17" wangid="0x22211122"/>
   <wangtile tileid="18" wangid="0x21212122"/>
   <wangtile tileid="19" wangid="0x21112122"/>
   <wangtile tileid="20" wangid="0x21211122"/>
   <wangtile tileid="21" wangid="0x21111122"/>
   <wangtile tileid="22" wangid="0x22222221"/>
   <wangtile tileid="23" wangid="0x21222221"/>
   <wangtile tileid="24" wangid="0x11222221"/>
   <wangtile tileid="25" wangid="0x22212221"/>
   <wangtile tileid="26" wangid="0x21212221"/>
   <wangtile tileid="27" wangid="0x11212221"/>
   <wangtile tileid="28" wangid="0x21112221"/>
   <wangtile tileid="29" wangid="0x11112221"/>
   <wangtile tileid="30" wangid="0x22222121"/>
   <wangtile tileid="31" wangid="0x22222111"/>
   <wangtile tileid="32" wangid="0x21222121"/>
   <wangtile tileid="33" wangid="0x11222121"/>
   <wangtile tileid="34" wangid="0x21222111"/>
   <wangtile tileid="35" wangid="0x11222111"/>
   <wangtile tileid="36" wangid="0x22212121"/>
   <wangtile tileid="37" wangid="0x22211121"/>
   <wangtile tileid="38" wangid="0x22212111"/>
   <wangtile tileid="39" wangid="0x22211111"/>
   <wangtile tileid="40" wangid="0x21212121"/>
   <wangtile tileid="41" wangid="0x11212121"/>
   <wangtile tileid="42" wangid="0x21112121"/>
   <wangtile tileid="43" wangid="0x11112121"/>
   <wangtile tileid="44" wangid="0x21211121"/>
   <wangtile tileid="45" wangid="0x11211121"/>
   <wangtile tileid="46" wangid="0x21111121"/>
   <wangtile tileid="47" wangid="0x11111121"/>
   <wangtile tileid="48" wangid="0x21212111"/>
   <wangtile tileid="49" wangid="0x11212111"/>
   <wangtile tileid="50" wangid="0x21112111"/>
   <wangtile tileid="51" wangid="0x11112111"/>
   <wangtile tileid="52" wangid="0x21211111"/>
   <wangtile tileid="53" wangid="0x11211111"/>
   <wangtile tileid="54" wangid="0x21111111"/>
   <wangtile tileid="55" wangid="0x11111111"/>
  </wangset>
 </wangsets>
</tileset>
