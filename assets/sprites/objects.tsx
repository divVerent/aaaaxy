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
  <image width="32" height="32" source="editorimgs/oneway_e.png"/>
 </tile>
 <tile id="13" type="OneWay">
  <properties>
   <property name="orientation" value="NW"/>
  </properties>
  <image width="32" height="32" source="editorimgs/oneway_n.png"/>
 </tile>
 <tile id="14" type="OneWay">
  <properties>
   <property name="orientation" value="WS"/>
  </properties>
  <image width="32" height="32" source="editorimgs/oneway_w.png"/>
 </tile>
 <tile id="15" type="OneWay">
  <properties>
   <property name="orientation" value="SE"/>
  </properties>
  <image width="32" height="32" source="editorimgs/oneway_s.png"/>
 </tile>
 <tile id="16" type="Switch">
  <image width="16" height="16" source="switch_off.png"/>
 </tile>
 <tile id="17" type="Riser">
  <image width="16" height="16" source="riser_small_idle.png"/>
 </tile>
 <tile id="18" type="SwitchableSprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value=""/>
   <property name="invert" type="bool" value="true"/>
  </properties>
  <image width="16" height="16" source="switchblock_off.png"/>
 </tile>
 <tile id="19" type="SwitchableSprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value=""/>
   <property name="invert" type="bool" value="false"/>
  </properties>
  <image width="16" height="16" source="switchblock_on.png"/>
 </tile>
 <tile id="20" type="TnihSign">
  <image width="32" height="32" source="tnihsign.png"/>
 </tile>
 <tile id="21" type="QuestionBlock">
  <image width="16" height="16" source="questionblock.png"/>
 </tile>
 <tile id="22" type="QuestionBlock">
  <image width="16" height="16" source="editorimgs/kaizoblock.png"/>
  <properties>
   <property name="kaizo" type="bool" value="true"/>
  </properties>
 </tile>
 <tile id="23" type="AppearBlock">
  <image width="16" height="16" source="appearblock.png"/>
 </tile>
 <tile id="24" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
  </properties>
  <image width="16" height="16" source="editorimgs/gradient_left_right.png"/>
 </tile>
 <tile id="25" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
  </properties>
  <image width="16" height="16" source="editorimgs/gradient_top_bottom.png"/>
 </tile>
 <tile id="26" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
  </properties>
  <image width="16" height="16" source="editorimgs/gradient_outside_inside.png"/>
 </tile>
 <tile id="27" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32s.png"/>
   <property name="orientation" value="ES"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_e.png"/>
 </tile>
 <tile id="28" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32d.png"/>
   <property name="orientation" value="NE"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_ne.png"/>
 </tile>
 <tile id="29" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32s.png"/>
   <property name="orientation" value="NE"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_n.png"/>
 </tile>
 <tile id="30" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32d.png"/>
   <property name="orientation" value="WN"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_nw.png"/>
 </tile>
 <tile id="31" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32s.png"/>
   <property name="orientation" value="WN"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_w.png"/>
 </tile>
 <tile id="32" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32d.png"/>
   <property name="orientation" value="SW"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_sw.png"/>
 </tile>
 <tile id="33" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32s.png"/>
   <property name="orientation" value="SW"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_s.png"/>
 </tile>
 <tile id="34" type="Sprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="image" value="arrow32d.png"/>
   <property name="orientation" value="ES"/>
  </properties>
  <image width="32" height="32" source="editorimgs/arrow32_se.png"/>
 </tile>
 <tile id="35" type="ForceField">
  <properties>
   <property name="orientation" value="ES"/>
  </properties>
  <image width="16" height="16" source="editorimgs/forcefield_v.png"/>
 </tile>
 <tile id="36" type="ForceField">
  <properties>
   <property name="orientation" value="SW"/>
  </properties>
  <image width="16" height="16" source="editorimgs/forcefield_h.png"/>
 </tile>
 <tile id="37" type="MovableSprite">
  <properties>
   <property name="image_dir" value=""/>
   <property name="delta" value="0 -64"/>
   <property name="solid" value="true"/>
  </properties>
  <image width="32" height="64" source="movingdoor.png"/>
 </tile>
</tileset>
